package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type ManufacturerPage struct {
	Manufacturer Manufacturer
	Cars         []CarModel
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryPage struct {
	Category Category
	Cars     []CarModel
}

type CarModel struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	ManufacturerID int            `json:"manufacturerId"`
	CategoryID     int            `json:"categoryId"`
	Year           int            `json:"year"`
	Specs          Specifications `json:"specifications"`
	Image          string         `json:"image"`
	Manufacturer   *Manufacturer  `json:"manufacturer,omitempty"`
	Category       *Category      `json:"category,omitempty"`
}

type Specifications struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// write the index.html page
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

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

	tmpl.Execute(w, cars)
}

func carHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/car.html"))

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

	catResp, catErr := http.Get(fmt.Sprintf("http://localhost:3000/api/categories/%d", car.CategoryID))
	if catErr != nil {
		log.Printf("category fetch failed for car %d: %v", car.ID, catErr)
	} else {
		defer catResp.Body.Close()

		var category Category

		if catDecodeErr := json.NewDecoder(catResp.Body).Decode(&category); catDecodeErr == nil {
			car.Category = &category
		} else {
			log.Printf("category decode failed for car %d: %v", car.ID, catDecodeErr)
		}
	}

	manuResp, manuErr := http.Get(fmt.Sprintf("http://localhost:3000/api/manufacturers/%d", car.ManufacturerID))
	if manuErr != nil {
		log.Printf("manufacturer fetch failed for car %d: %v", car.ID, manuErr)
	} else {
		defer manuResp.Body.Close()

		var manufacturer Manufacturer

		if manuDecodeErr := json.NewDecoder(manuResp.Body).Decode(&manufacturer); manuDecodeErr == nil {
			car.Manufacturer = &manufacturer
		} else {
			log.Printf("manufacturer decode failed for car %d: %v", car.ID, manuDecodeErr)
		}
	}
	tmpl.Execute(w, car)

}

func categoryHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/category.html"))

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

	catResp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/categories/%s", id))
	if err != nil {
		http.Error(w, "fetching category failed", http.StatusInternalServerError)
		return
	}
	defer catResp.Body.Close()

	var category Category

	if err := json.NewDecoder(catResp.Body).Decode(&category); err != nil {
		http.Error(w, "data error decoding category", http.StatusInternalServerError)
		return
	}

	var filtered []CarModel

	for _, c := range cars {
		if c.CategoryID == category.ID {
			filtered = append(filtered, c)
		}
	}

	tmpl.Execute(w, CategoryPage{Category: category, Cars: filtered})

}

func manufacturerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/manufacturer.html"))

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

	manuResp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/manufacturers/%s", id))
	if err != nil {
		http.Error(w, "fetching manufacturer failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var manufacturer Manufacturer

	if err := json.NewDecoder(manuResp.Body).Decode(&manufacturer); err != nil {
		http.Error(w, "data error decoding  manufacturer", http.StatusInternalServerError)
		return
	}

	var filtered []CarModel

	for _, c := range cars {
		if c.ManufacturerID == manufacturer.ID {
			filtered = append(filtered, c)
		}
	}

	tmpl.Execute(w, ManufacturerPage{Manufacturer: manufacturer, Cars: filtered})
}
