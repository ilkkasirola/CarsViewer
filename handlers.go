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
	// write the index.html page
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
	// get the file handle and check for errors
	file, err := os.Open("api/data.json")
	if err != nil {
		http.Error(w, "json file not found", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var d Data
	// read the json data and convert it to our go struct
	err = json.NewDecoder(file).Decode(&d)
	if err != nil {
		http.Error(w, "data error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, d.CarModels)
}
