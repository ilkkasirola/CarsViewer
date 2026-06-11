package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", homeHandler)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /car/{id}", carHandler)
	mux.HandleFunc("GET /category/{id}", categoryHandler)
	mux.HandleFunc("GET /manufacturer/{id}", manufacturerHandler)
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
