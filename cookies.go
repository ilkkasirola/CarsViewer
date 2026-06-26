package main

import (
	"net/http"
	"net/url"
	"slices"
	"sort"
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

func addToCompare(w http.ResponseWriter, r *http.Request, currentID int) {

	var compareID = getCompareIDs(r)

	if !slices.Contains(compareID, currentID) && len(compareID) < 2 {
		compareID = append(compareID, currentID)
	}

	var parts []string
	for _, id := range compareID {
		parts = append(parts, strconv.Itoa(id))
	}

	newValue := strings.Join(parts, ",")
	http.SetCookie(w, &http.Cookie{
		Name:  "compareID",
		Value: newValue,
		Path:  "/",
	})

}

func removeFromCompare(w http.ResponseWriter, r *http.Request, currentID int) {
	var compareID = getCompareIDs(r)

	var filtered []int
	for _, id := range compareID {
		if id != currentID {
			filtered = append(filtered, id)
		}
	}

	var parts []string
	for _, id := range filtered {
		parts = append(parts, strconv.Itoa(id))
	}

	newValue := strings.Join(parts, ",")
	http.SetCookie(w, &http.Cookie{
		Name:  "compareID",
		Value: newValue,
		Path:  "/",
	})
}

func getCompareIDs(r *http.Request) []int {
	var ids []int
	if car, err := r.Cookie("compareID"); err == nil && car.Value != "" {
		for _, s := range strings.Split(car.Value, ",") {
			if s == "" {
				continue
			}
			if id, err := strconv.Atoi(s); err == nil {
				ids = append(ids, id)
			}
		}
	}
	return ids

}

func getRecentlyViewed(w http.ResponseWriter, r *http.Request, currentID int, allCars []CarModel, maxN int) ([]CarModel, error) {

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
	return recentViewed, nil
}

func giveRecommendations(recents []CarModel, allCars []CarModel, maxN int) ([]CarModel, error) {
	if len(recents) == 0 || maxN <= 0 {
		return []CarModel{}, nil
	}
	seen := map[int]struct{}{}
	countries := map[string]struct{}{}
	categories := map[int]struct{}{}
	hpCounts := map[int]int{}

	for _, r := range recents {
		seen[r.ID] = struct{}{}
		if r.Manufacturer != nil && r.Manufacturer.Country != "" {
			countries[r.Manufacturer.Country] = struct{}{}
		}
		if r.CategoryID != 0 {
			categories[r.CategoryID] = struct{}{}
		}
		if r.Specs.Horsepower > 0 {
			hpCounts[r.Specs.Horsepower]++
		}
	}
	score := func(c CarModel) int {
		s := 0
		if c.Manufacturer != nil {
			if _, ok := countries[c.Manufacturer.Country]; ok && c.Manufacturer.Country != "" {
				s += 100
			}
		}
		if c.CategoryID != 0 {
			if _, ok := categories[c.CategoryID]; ok {
				s += 3
			}
		}
		if c.Specs.Horsepower > 0 {
			s += hpCounts[c.Specs.Horsepower] * 2
		}
		return s
	}
	cands := make([]CarModel, 0, len(allCars))
	for _, c := range allCars {
		if _, ok := seen[c.ID]; !ok {
			cands = append(cands, c)
		}
	}

	sort.Slice(cands, func(i, j int) bool {
		return score(cands[i]) > score(cands[j])
	})

	if len(cands) > maxN {
		cands = cands[:maxN]
	}
	return cands, nil
}
