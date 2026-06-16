package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func loadAllCars() ([]CarModel, error) {
	file, err := os.ReadFile("data/cars.json")
	if err != nil {
		return nil, err
	}
	var cars []CarModel
	if err := json.Unmarshal(file, &cars); err != nil {
		return nil, err
	}
	return cars, nil

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
	allCars, err := loadAllCars()
	if err != nil {
		return nil, err
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
