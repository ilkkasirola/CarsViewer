package main

import (
	"encoding/json"
	"fmt"
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
	// write the index.html page
	w.Header().Set("Content-Type", "text/html")

	// err := template.Must(template.ParseFiles("templates/index.html", "templates/topnav.html"))

	// get the response from api endpoint
	resp, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		http.Error(w, "fetching failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var cars []CarModel
	err = json.NewDecoder(resp.Body).Decode(&cars)
	// read the json data and convert it to our go struct
	if err != nil {
		http.Error(w, "data error", http.StatusInternalServerError)
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

	nav := fetchNav()
	nav.SelectedCategories = categoryIDs
	nav.SelectedManufacturers = manufacturerIDs

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
	// tmpl := template.Must(template.ParseFiles("templates/car.html", "templates/topnav.html"))

	id := strings.TrimPrefix(r.URL.Path, "/car/")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	resp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/models/%s", id))
	if err != nil {
		http.Error(w, "fetching failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var car CarModel
	err = json.NewDecoder(resp.Body).Decode(&car)
	if err != nil {
		http.Error(w, "data error decoding car", http.StatusInternalServerError)
		return
	}

	nav := fetchNav()
	//using nav to get categories and manufacturers reducing api calls
	for _, c := range nav.Categories {
		if c.ID == car.CategoryID {
			car.Category = &c
			break
		}
	}

	for _, m := range nav.Manufacturers {
		if m.ID == car.ManufacturerID {
			car.Manufacturer = &m
			break
		}
	}

	recents, err := getRecentlyViewed(w, r, car.ID, 5)
	if err != nil {
		log.Printf("getRecentlu error: %v", err)
		http.Error(w, "cannot get recently viewed", http.StatusInternalServerError)
		return
	}
	// using referer for "back" button to keep filters when going back to the home page
	referer := r.Referer()
	if referer == "" {
		referer = "/"
	}

	data := CarPage{
		Nav:            nav,
		Car:            car,
		RecentlyViewed: recents,
		BackURL:        referer,
	}

	if err := templates.ExecuteTemplate(w, "car.html", data); err != nil {
		log.Printf("template exxecute error: %v", err)
	}
}

func fetchNav() Nav {
	// helper function fetching manufacturers and categories for navigation bar
	var nav Nav
	if manuResp, manuErr := http.Get("http://localhost:3000/api/manufacturers"); manuErr == nil {
		defer manuResp.Body.Close()
		json.NewDecoder(manuResp.Body).Decode(&nav.Manufacturers)
	} else {
		log.Printf("manufacturers fetch failed: %v", manuErr)
	}
	if catResp, catErr := http.Get("http://localhost:3000/api/categories"); catErr == nil {
		defer catResp.Body.Close()
		json.NewDecoder(catResp.Body).Decode(&nav.Categories)
	} else {
		log.Printf("categories fetch failed %v", catErr)
	}
	return nav
}

// convert ids to ints
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
