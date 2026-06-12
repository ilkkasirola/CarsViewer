package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const cookieName = "visitorID"

func newVisitorID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + "-" + strconv.Itoa(rand.Intn(1000000))
}
func getOrSetVisitorCookie(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie(cookieName)
	if err == nil && c.Value != "" {
		return c.Value
	}
	id := newVisitorID()
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    id,
		MaxAge:   3600,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		// Secure:   true,
	}
	http.SetCookie(w, &cookie)
	return id
}
