package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"counter/backend"

	"github.com/gorilla/mux"
)

func main() {
	// Print the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	log.Println("Current working directory:", dir)

	// Initialize the database
	err = backend.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer backend.CloseDB()

	router := mux.NewRouter()

	// Serve static files from the "public" directory with correct MIME types
	fs := http.Dir("./public")
	fileServer := http.FileServer(fs)
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.URL.Path)
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".ttf":
			w.Header().Set("Content-Type", "font/ttf")
		}
		fileServer.ServeHTTP(w, r)
	})))

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
