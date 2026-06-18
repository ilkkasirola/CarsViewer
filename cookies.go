package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func saveFiltersCookie(w http.ResponseWriter, query string) {
	http.SetCookie(w, &http.Cookie{
		Name: "filters",
		// had some problems with '=' character. have to use QueryEscape
		Value: url.QueryEscape(query),
		Path:  "/",
	})
}

func getFilterBackURL(r *http.Request) string {
	if c, err := r.Cookie("filters"); err == nil && c.Value != "" {
		if q, err := url.QueryUnescape(c.Value); err == nil {
			return "/?" + q
		}
	}
	return "/"
}

func getRecentlyViewed(w http.ResponseWriter, r *http.Request, currentID int, maxN int) ([]CarModel, error) {

	var viewed []int
	if car, err := r.Cookie("view_history"); err == nil && car.Value != "" {
		for _, s := range strings.Split(car.Value, ",") {
			if s == "" {
				continue
			}
			if id, err := strconv.Atoi(s); err == nil {
				viewed = append(viewed, id)
			}
		}
	}

	if currentID != 0 {
		if len(viewed) == 0 || viewed[len(viewed)-1] != currentID {
			viewed = append(viewed, currentID)
			if len(viewed) > 10 {
				viewed = viewed[len(viewed)-10:]
			}
		}
	}

	var parts []string
	for _, id := range viewed {
		parts = append(parts, strconv.Itoa(id))
	}

	newValue := strings.Join(parts, ",")
	http.SetCookie(w, &http.Cookie{
		Name:  "view_history",
		Value: newValue,
		Path:  "/",
	})
	if len(viewed) == 0 || (len(viewed) == 1 && viewed[0] == currentID) {
		return nil, nil
	}

	resp, err := http.Get("http://localhost:3000/api/models")
	if err != nil {
		log.Printf("failed to fetch models from API: %v", err)
		return nil, nil
	}
	defer resp.Body.Close()

	var allCars []CarModel
	if err := json.NewDecoder(resp.Body).Decode(&allCars); err != nil {
		log.Printf("failed to decode models JSON: %v", err)
		return nil, nil
	}

	byID := map[int]CarModel{}
	for _, car := range allCars {
		byID[car.ID] = car
	}

	var recentViewed []CarModel
	for i := len(viewed) - 1; i >= 0 && len(recentViewed) < maxN; i-- {
		id := viewed[i]
		if id == currentID {
			continue
		}
		if car, ok := byID[id]; ok {
			dupe := false
			for _, recentCar := range recentViewed {
				if recentCar.ID == car.ID {
					dupe = true
					break
				}
			}
			if !dupe {
				recentViewed = append(recentViewed, car)
			}
		}
	}

	// if len(recentViewed) < limit {
	// 	for _, c := range allCars {
	// 		if len(recentViewed) >= limit
	// 	}
	// }
	return recentViewed, nil
}
