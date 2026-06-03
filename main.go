package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", homeHandler)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("GET /api/img/", http.StripPrefix("/api/img/", http.FileServer(http.Dir("api/img"))))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
