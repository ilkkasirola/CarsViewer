package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func fetchHomeData() (Nav, []CarModel, error) {

	manuCh := make(chan []Manufacturer, 1)
	catCh := make(chan []Category, 1)
	carsCh := make(chan []CarModel, 1)
	errCh := make(chan error, 3)

	go func() {
		var v []Manufacturer
		if resp, err := http.Get("http://localhost:3000/api/manufacturers"); err != nil {
			errCh <- err
		} else {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&v)
			manuCh <- v
		}
	}()

	go func() {
		var v []Category
		if resp, err := http.Get("http://localhost:3000/api/categories"); err != nil {
			errCh <- err
		} else {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&v)
			catCh <- v
		}
	}()

	go func() {
		var v []CarModel
		if resp, err := http.Get("http://localhost:3000/api/models"); err != nil {
			errCh <- err
		} else {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&v)
			carsCh <- v
		}
	}()

	var nav Nav
	var cars []CarModel
	for range 3 {
		select {
		case v := <-manuCh:
			nav.Manufacturers = v
		case v := <-catCh:
			nav.Categories = v
		case v := <-carsCh:
			cars = v
		case err := <-errCh:
			return Nav{}, nil, err

		}
	}
	return nav, cars, nil
}

func fetchCarData(carID string) (Lookup, CarModel, error) {
	manuCh := make(chan []Manufacturer, 1)
	catCh := make(chan []Category, 1)
	carCh := make(chan CarModel, 1)
	errCh := make(chan error, 3)

	go func() {
		var v []Manufacturer
		if resp, err := http.Get("http://localhost:3000/api/manufacturers"); err != nil {
			errCh <- err
		} else {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&v)
			manuCh <- v
		}
	}()

	go func() {
		var v []Category
		if resp, err := http.Get("http://localhost:3000/api/categories"); err != nil {
			errCh <- err
		} else {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&v)
			catCh <- v
		}
	}()

	go func() {
		var v CarModel
		if resp, err := http.Get(fmt.Sprintf("http://localhost:3000/api/models/%s", carID)); err != nil {
			errCh <- err
		} else {
			defer resp.Body.Close()
			json.NewDecoder(resp.Body).Decode(&v)
			carCh <- v
		}
	}()
	var lookup Lookup
	var car CarModel
	for range 3 {
		select {
		case v := <-manuCh:
			lookup.Manufacturers = v
		case v := <-catCh:
			lookup.Categories = v
		case v := <-carCh:
			car = v
		case err := <-errCh:
			return Lookup{}, CarModel{}, err

		}
	}
	return lookup, car, nil

}
