package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

var templates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/car.html",
	"templates/topnav.html",
))

func homeHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	nav, cars, err := fetchHomeData()
	if err != nil {
		http.Error(w, "fetching failed", http.StatusInternalServerError)
		return
	}

	categoryIDs, err := parseIDs(r.URL.Query()["category"])
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	manufacturerIDs, err := parseIDs(r.URL.Query()["manufacturer"])
	if err != nil {
		http.Error(w, "invalid manufacturer id", http.StatusBadRequest)
		return
	}

	nav.SelectedCategories = categoryIDs
	nav.SelectedManufacturers = manufacturerIDs

	//save filters from main page
	if r.URL.Path == "/" {
		saveFiltersCookie(w, r.URL.RawQuery)
	}

	var filtered []CarModel
	for _, c := range cars {
		if len(categoryIDs) > 0 && !slices.Contains(categoryIDs, c.CategoryID) {
			continue
		}
		if len(manufacturerIDs) > 0 && !slices.Contains(manufacturerIDs, c.ManufacturerID) {
			continue
		}
		filtered = append(filtered, c)
	}

	err = templates.ExecuteTemplate(w, "index.html", HomePage{Nav: nav, Cars: filtered})
	if err != nil {
		log.Printf("templates execute  error: %v", err)
	}

}

func carHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	id := strings.TrimPrefix(r.URL.Path, "/car/")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	lookup, car, err := fetchCarData(id)
	if err != nil {
		http.Error(w, "fetching failed", http.StatusInternalServerError)
		return
	}

	for _, c := range lookup.Categories {
		if c.ID == car.CategoryID {
			car.Category = &c
			break
		}
	}

	for _, m := range lookup.Manufacturers {
		if m.ID == car.ManufacturerID {
			car.Manufacturer = &m
			break
		}
	}

	recents, err := getRecentlyViewed(w, r, car.ID, 5)
	if err != nil {
		log.Printf("getRecently error: %v", err)
		http.Error(w, "cannot get recently viewed", http.StatusInternalServerError)
		return
	}
	resp, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		log.Printf("failed to fetch models from API: %v", err)
		http.Error(w, "cannot get models", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var allCars []CarModel
	if err := json.NewDecoder(resp.Body).Decode(&allCars); err != nil {
		log.Printf("failed to decode models JSON: %v", err)
		http.Error(w, "cannot deode models", http.StatusInternalServerError)
	}
	recs, err := giveRecommendations(recents, allCars, 10)
	if err != nil {
		log.Printf("reccommendations error: %v", err)
		recs = []CarModel{}
	}

	backURL := getFilterBackURL(r)

	data := CarPage{
		Lookup:          lookup,
		Car:             car,
		RecentlyViewed:  recents,
		Recommendations: recs,
		BackURL:         backURL,
	}

	if err := templates.ExecuteTemplate(w, "car.html", data); err != nil {
		log.Printf("template execute error: %v", err)
	}
}

func parseIDs(values []string) ([]int, error) {
	ids := []int{}
	for _, v := range values {
		id, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// check filters to render for checkboxes
func (n Nav) IsManufacturerSelected(id int) bool {
	return slices.Contains(n.SelectedManufacturers, id)
}

func (n Nav) IsCategorySelected(id int) bool {
	return slices.Contains(n.SelectedCategories, id)
}
