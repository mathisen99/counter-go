package main

import (
	"log"
	"net/http"

	"counter/backend" // Replace with your actual import path
)

func main() {
	// Initialize the database
	err := backend.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer backend.CloseDB()

	// Serve static files
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// Register handlers
	http.HandleFunc("/", backend.HomeHandler)
	http.HandleFunc("/admin", backend.AdminHandler)
	http.HandleFunc("/admin", backend.AdminPostHandler)
	http.HandleFunc("/update", backend.UpdateHandler)

	// Start the server
	log.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
