package backend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/nfnt/resize"
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
		"TimeCounterDisplay": settings.TimeCounterDisplay,
		"FontSize":           settings.FontSize,
		"LogoText":           settings.LogoText,
		"LogoFontSize":       settings.LogoFontSize,
		"Fields":             fields,
	})
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

// AdminHandler handles the admin page
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
		Images:   images,
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
	fontSize, err := strconv.Atoi(r.FormValue("font_size"))
	if err != nil {
		http.Error(w, "Invalid font size value", http.StatusBadRequest)
		return
	}
	logoText := r.FormValue("logo_text")
	logoFontSize, err := strconv.Atoi(r.FormValue("logo_font_size"))
	if err != nil {
		http.Error(w, "Invalid logo font size value", http.StatusBadRequest)
		return
	}

	// Update settings with form values
	settings.Start = start
	settings.Time = time.Now().Unix()
	settings.Speed = speed
	settings.ViewSpeed = viewSpeed
	settings.TimeCounterDisplay = timeCounterDisplay
	settings.FontSize = fontSize
	settings.LogoText = logoText
	settings.LogoFontSize = logoFontSize
	saveSettings()

	// Parse form values
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Loop over form values field_text and show_field and save them to the database
	for i := 1; i <= 5; i++ {
		fieldText := r.FormValue(fmt.Sprintf("field_text%d", i))
		showField := r.FormValue(fmt.Sprintf("show_field%d", i))

		if fieldText != "" {
			if showField == "" {
				showField = "off"
			}

			// Use custom function to save the field to the database
			saveField(i, fieldText, showField)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	filePath := filepath.Join("./images", header.Filename)
	out, err := os.Create(filePath)
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

	// Generate thumbnail
	err = generateThumbnail(filePath, header.Filename)
	if err != nil {
		log.Printf("Error generating thumbnail: %v", err)
		http.Error(w, "Failed to generate thumbnail", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func generateThumbnail(filePath, fileName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	var img image.Image
	if filepath.Ext(fileName) == ".png" {
		img, err = png.Decode(file)
	} else if filepath.Ext(fileName) == ".jpg" || filepath.Ext(fileName) == ".jpeg" {
		img, err = jpeg.Decode(file)
	} else {
		return fmt.Errorf("unsupported image format")
	}
	if err != nil {
		return err
	}

	// Resize the image to a thumbnail
	thumb := resize.Resize(100, 0, img, resize.Lanczos3)

	// Create thumbnail file
	thumbPath := filepath.Join("./images", "thumb_"+fileName)
	out, err := os.Create(thumbPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Encode the thumbnail to disk
	if filepath.Ext(fileName) == ".png" {
		err = png.Encode(out, thumb)
	} else {
		err = jpeg.Encode(out, thumb, nil)
	}
	return err
}
