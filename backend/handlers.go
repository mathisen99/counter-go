package backend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	tmplCache = map[string]*template.Template{}
	tmplMutex sync.Mutex
)

func loadTemplate(name string) (*template.Template, error) {
	tmplMutex.Lock()
	defer tmplMutex.Unlock()

	if tmpl, ok := tmplCache[name]; ok {
		return tmpl, nil
	}

	tmplPath := filepath.Join("templates", name)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template %s: %w", name, err)
	}

	tmplCache[name] = tmpl
	return tmpl, nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Reload settings
	LoadSettings()

	// Get current time as seconds since Unix epoch
	curTime := time.Now().Unix()

	// Calculate time passed in seconds
	timePassed := curTime - settings.Time

	// Calculate new start value
	newStart := settings.Start + (int(timePassed)*1000)/settings.Speed

	// Load fields from the database
	fields, err := LoadFields()
	if err != nil {
		log.Printf("Error loading fields from database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := loadTemplate("index.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"Start":              newStart,
		"Speed":              settings.Speed,
		"ViewSpeed":          settings.ViewSpeed,
		"TimeCounterDisplay": settings.TimeCounterDisplay, // Pass the new setting
		"Fields":             fields,
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	// Reload settings
	LoadSettings()

	tmpl, err := loadTemplate("admin.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Current time as seconds since Unix epoch
	curTime := time.Now().Unix()

	// Calculate time passed in seconds
	timePassed := curTime - settings.Time

	// Calculate new start value for display purposes only
	newStart := settings.Start + (int(timePassed)*1000)/settings.Speed

	// Create a copy of settings to display without altering the original settings
	displaySettings := settings
	displaySettings.Start = newStart

	// Load fields from the database
	fields, err := LoadFields()
	if err != nil {
		log.Printf("Error loading fields from database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Ensure there are at least 5 fields
	for len(fields) < 5 {
		fields = append(fields, Field{
			ID:        len(fields) + 1,
			FieldText: "",
			ShowField: "off",
		})
	}

	// Load images from the images directory
	var images []string
	files, err := os.ReadDir("./images")
	if err == nil {
		for _, file := range files {
			if !file.IsDir() {
				images = append(images, file.Name())
			}
		}
	}

	// Create a DisplayData struct to pass to the template
	displayData := DisplayData{
		Settings: displaySettings,
		Fields:   fields,
		Images:   images, // Pass images to the template
	}

	// Populate the fields with the current values
	err = tmpl.Execute(w, displayData)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func AdminPostHandler(w http.ResponseWriter, r *http.Request) {
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil {
		http.Error(w, "Invalid start value", http.StatusBadRequest)
		return
	}
	speed, err := strconv.Atoi(r.FormValue("speed"))
	if err != nil {
		http.Error(w, "Invalid speed value", http.StatusBadRequest)
		return
	}
	viewSpeed, err := strconv.Atoi(r.FormValue("view_speed"))
	if err != nil {
		http.Error(w, "Invalid view speed value", http.StatusBadRequest)
		return
	}
	timeCounterDisplay, err := strconv.Atoi(r.FormValue("time_counter_display"))
	if err != nil {
		http.Error(w, "Invalid time counter display value", http.StatusBadRequest)
		return
	}

	settings.Start = start
	settings.Time = time.Now().Unix()
	settings.Speed = speed
	settings.ViewSpeed = viewSpeed
	settings.TimeCounterDisplay = timeCounterDisplay
	saveSettings()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	for i := 1; i <= 5; i++ {
		fieldText := r.FormValue(fmt.Sprintf("field_text%d", i))
		showField := r.FormValue(fmt.Sprintf("show_field%d", i))
		imgSize := r.FormValue(fmt.Sprintf("img_size%d", i))

		if fieldText != "" {
			if showField == "" {
				showField = "off"
			}

			saveField(i, fieldText, showField, imgSize)
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil {
		http.Error(w, "Invalid start value", http.StatusBadRequest)
		return
	}
	timeVal, err := strconv.Atoi(r.FormValue("time"))
	if err != nil {
		http.Error(w, "Invalid time value", http.StatusBadRequest)
		return
	}

	settings.Start = start
	settings.Time = int64(timeVal)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// New handler for image upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error getting uploaded file: %v", err)
		http.Error(w, "Failed to upload file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the images directory if it doesn't exist
	os.MkdirAll("./images", os.ModePerm)

	// Save the file to the images directory
	out, err := os.Create(filepath.Join("./images", header.Filename))
	if err != nil {
		log.Printf("Error saving uploaded file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Printf("Error copying uploaded file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
