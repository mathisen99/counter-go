package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"counter/backend"
)

func main() {
	// Initialize the database
	err := backend.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer backend.CloseDB()

	router := mux.NewRouter()

	// Serve static files from the "public" directory
	fs := http.FileServer(http.Dir("public"))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	// Define routes
	router.HandleFunc("/", backend.HomeHandler).Methods("GET")
	router.HandleFunc("/admin", backend.AdminHandler).Methods("GET")
	router.HandleFunc("/admin", backend.AdminPostHandler).Methods("POST")
	router.HandleFunc("/update", backend.UpdateHandler).Methods("POST")

	// Load settings on startup
	backend.LoadSettings()

	log.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
