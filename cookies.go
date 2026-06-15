package main

import (
	"net/http"
	"strconv"
	"strings"
)

func recommendHandler(w http.ResponseWriter, r *http.Request) {
	currentItem := r.URL.Query().Get("id")
	if currentItem == "" {
		currentItem = "home"
	}
	var viewHistoryIDs []int
	cookie, err := r.Cookie("view_history")
	if err == nil && cookie.Value != "" {
		parts := strings.Split(cookie.Value, ",")
		for _, p := range parts {
			if p == "" {
				continue
			}
			id, err := strconv.Atoi(p)
			if err == nil {
				viewHistoryIDs = append(viewHistoryIDs, id)
			}
		}
	}
	if len(viewHistoryIDs) == 0 || strconv.Itoa(viewHistoryIDs[len(viewHistoryIDs)-1]) != currentItem {
		if id, err := strconv.Atoi(currentItem); err == nil {
			viewHistoryIDs = append(viewHistoryIDs, id)
		}
		if len(viewHistoryIDs) > 10 {
			viewHistoryIDs = viewHistoryIDs[len(viewHistoryIDs)-10:]
		}
	}
	newValueParts := make([]string, 0, len(viewHistoryIDs))
	for _, id := range viewHistoryIDs {
		newValueParts = append(newValueParts, strconv.Itoa(id))
	}
	newValue := strings.Join(newValueParts, ",")

	http.SetCookie(w, &http.Cookie{
		Name:     "view_history",
		Value:    newValue,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}
