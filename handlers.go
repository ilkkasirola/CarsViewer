package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CarModel struct {
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	ManufacturerID int            `json:"manufacturerId"`
	CategoryID     int            `json:"categoryId"`
	Year           int            `json:"year"`
	Specs          Specifications `json:"specifications"`
	Image          string         `json:"image"`
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
	tmpl := template.Must(template.ParseFiles("index.html"))

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
	tmpl := template.Must(template.ParseFiles("car.html"))

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

	tmpl.Execute(w, car)
}
