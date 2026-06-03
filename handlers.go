package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

type Data struct {
	Manufacturers []Manufacturer   `json:"manufacturers"`
	CarModels     []CarModel       `json:"carModels"`
	Categories    []Category       `json:"categories"`
	Specs         []Specifications `json:"specifications"`
}

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
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadRequest)

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
	file, err := os.Open("api/data.json")
	if err != nil {
		http.Error(w, "data not found", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var d Data

	err = json.NewDecoder(file).Decode(&d)
	if err != nil {
		http.Error(w, "bad data", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, d.CarModels)
}
