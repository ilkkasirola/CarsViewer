package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/car/", carHandler)
	mux.HandleFunc("POST /compare/add/", compareAddHandler)
	mux.HandleFunc("POST /compare/remove/", compareRemoveHandler)
	mux.HandleFunc("/compare", comparePageHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
