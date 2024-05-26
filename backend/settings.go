package backend

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Settings struct {
	Start              int    `json:"start"`
	Speed              int    `json:"speed"`
	Time               int64  `json:"time"`
	ViewSpeed          int    `json:"view_speed"`
	TimeCounterDisplay int    `json:"time_counter_display"`
	FontSize           int    `json:"font_size"`
	LogoText           string `json:"logo_text"`
	LogoFontSize       int    `json:"logo_font_size"` // New field
}

type Field struct {
	ID        int
	FieldText string
	ShowField string
}

type DisplayData struct {
	Settings   Settings
	Fields     []Field
	ShowFields []string
	Images     []string
}

var settings Settings
var db *sql.DB

func InitDB() error {
	var err error
	// Check if the database file exists
	if _, err := os.Stat("./settings.db"); os.IsNotExist(err) {
		// Database file does not exist, create it
		log.Println("Creating new database file...")
		file, err := os.Create("./settings.db")
		if err != nil {
			return err
		}
		file.Close()
	}

	// Open the database
	db, err = sql.Open("sqlite3", "./settings.db")
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY,
		start INTEGER,
		speed INTEGER,
		time INTEGER,
		view_speed INTEGER,
		time_counter_display INTEGER DEFAULT 5000,
		font_size INTEGER DEFAULT 24,
		logo_text TEXT DEFAULT 'Pengar',
		logo_font_size INTEGER DEFAULT 36
	)`)
	if err != nil {
		return err
	}

	// Create table for fields
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS fields (
		id INTEGER PRIMARY KEY,
		field_text TEXT,
		show_field TEXT
	)`)
	if err != nil {
		return err
	}

	// Insert default settings if they do not exist
	_, err = db.Exec(`INSERT INTO settings (id, start, speed, time, view_speed, time_counter_display, font_size, logo_text, logo_font_size)
		VALUES (1, 0, 1000, 0, 1000, 5000, 24, 'Pengar', 36)
		ON CONFLICT(id) DO NOTHING`)
	if err != nil {
		return err
	}

	return nil
}

func LoadSettings() {
	row := db.QueryRow("SELECT start, speed, time, view_speed, time_counter_display, font_size, logo_text, logo_font_size FROM settings WHERE id = 1")
	var start, speed, viewSpeed, timeCounterDisplay, fontSize, logoFontSize int
	var logoText string
	var time int64
	err := row.Scan(&start, &speed, &time, &viewSpeed, &timeCounterDisplay, &fontSize, &logoText, &logoFontSize)
	if err != nil {
		log.Println("Settings not found in database, using default settings")
		settings = Settings{Start: 0, Speed: 1000, ViewSpeed: 1000, TimeCounterDisplay: 5000, FontSize: 24, LogoText: "Pengar", LogoFontSize: 36}
		return
	}
	settings = Settings{Start: start, Speed: speed, Time: time, ViewSpeed: viewSpeed, TimeCounterDisplay: timeCounterDisplay, FontSize: fontSize, LogoText: logoText, LogoFontSize: logoFontSize}
}

func saveSettings() {
	_, err := db.Exec("REPLACE INTO settings (id, start, speed, time, view_speed, time_counter_display, font_size, logo_text, logo_font_size) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, settings.Start, settings.Speed, settings.Time, settings.ViewSpeed, settings.TimeCounterDisplay, settings.FontSize, settings.LogoText, settings.LogoFontSize)
	if err != nil {
		log.Fatalf("Error saving settings to database: %v", err)
	}
}

func saveField(id int, fieldText, showField string) {
	_, err := db.Exec("REPLACE INTO fields (id, field_text, show_field) VALUES (?, ?, ?)", id, fieldText, showField)
	if err != nil {
		log.Fatalf("Error saving field to database: %v", err)
	}
}

func LoadFields() ([]Field, error) {
	rows, err := db.Query("SELECT id, field_text, show_field FROM fields")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []Field
	for rows.Next() {
		var field Field
		err := rows.Scan(&field.ID, &field.FieldText, &field.ShowField)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func CloseDB() {
	db.Close()
}
