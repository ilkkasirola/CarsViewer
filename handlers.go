package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// write the index.html page
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/topnav.html"))

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
	tmpl.Execute(w, HomePage{Nav: fetchNav(), Cars: cars})
}

func carHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/car.html", "templates/topnav.html"))

	id := r.PathValue("id")

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

	tmpl.Execute(w, CarPage{Nav: nav, Car: car})

}

func categoryHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/category.html", "templates/topnav.html"))

	id := r.PathValue("id")

	resp, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		http.Error(w, "fetching models failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var cars []CarModel

	err = json.NewDecoder(resp.Body).Decode(&cars)
	if err != nil {
		http.Error(w, "data error", http.StatusInternalServerError)
		return
	}

	nav := fetchNav()

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "error converting category id", http.StatusBadRequest)
		return
	}

	var category Category

	for _, c := range nav.Categories {
		if c.ID == categoryID {
			category = c
			break
		}
	}

	var filtered []CarModel

	for _, c := range cars {
		if c.CategoryID == category.ID {
			filtered = append(filtered, c)
		}
	}

	tmpl.Execute(w, CategoryPage{Nav: nav, Category: category, Cars: filtered})

}

func manufacturerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/manufacturer.html", "templates/topnav.html"))

	id := r.PathValue("id")

	resp, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		http.Error(w, "fetching models failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var cars []CarModel

	err = json.NewDecoder(resp.Body).Decode(&cars)
	if err != nil {
		http.Error(w, "data error", http.StatusInternalServerError)
		return
	}

	nav := fetchNav()

	manufacturerID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "error converting manufacturer id", http.StatusBadRequest)
		return
	}

	var manufacturer Manufacturer

	for _, m := range nav.Manufacturers {
		if m.ID == manufacturerID {
			manufacturer = m
			break
		}
	}

	var filtered []CarModel

	for _, c := range cars {
		if c.ManufacturerID == manufacturer.ID {
			filtered = append(filtered, c)
		}
	}

	tmpl.Execute(w, ManufacturerPage{Nav: nav, Manufacturer: manufacturer, Cars: filtered})
}
