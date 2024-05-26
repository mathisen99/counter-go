package backend

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Settings struct {
	Start              int   `json:"start"`
	Speed              int   `json:"speed"`
	Time               int64 `json:"time"`
	ViewSpeed          int   `json:"view_speed"`
	TimeCounterDisplay int   `json:"time_counter_display"`
}

type Field struct {
	ID        int
	FieldText string
	ShowField string
	ImgSize   string
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
	if _, err := os.Stat("./settings.db"); os.IsNotExist(err) {
		log.Println("Creating new database file...")
		file, err := os.Create("./settings.db")
		if err != nil {
			return err
		}
		file.Close()
	}

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
		time_counter_display INTEGER DEFAULT 5000
	)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS fields (
		id INTEGER PRIMARY KEY,
		field_text TEXT,
		show_field TEXT,
		img_size TEXT
	)`)

	if err != nil {
		return err
	}

	err = insertDefaultSettings()
	if err != nil {
		return err
	}

	return nil
}

func insertDefaultSettings() error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM settings WHERE id = 1").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = db.Exec("INSERT INTO settings (id, start, speed, time, view_speed, time_counter_display) VALUES (?, ?, ?, ?, ?, ?)",
			1, 0, 1000, 0, 1000, 5000)
		if err != nil {
			return err
		}
		log.Println("Inserted default settings into database")
	}

	return nil
}

func LoadSettings() {
	row := db.QueryRow("SELECT start, speed, time, view_speed, time_counter_display FROM settings WHERE id = 1")
	var start, speed, viewSpeed, timeCounterDisplay int
	var time int64
	err := row.Scan(&start, &speed, &time, &viewSpeed, &timeCounterDisplay)
	if err != nil {
		log.Println("Settings not found in database, using default settings")
		settings = Settings{Start: 0, Speed: 1000, ViewSpeed: 1000, TimeCounterDisplay: 5000}
		return
	}
	settings = Settings{Start: start, Speed: speed, Time: time, ViewSpeed: viewSpeed, TimeCounterDisplay: timeCounterDisplay}
}

func saveSettings() {
	_, err := db.Exec("REPLACE INTO settings (id, start, speed, time, view_speed, time_counter_display) VALUES (?, ?, ?, ?, ?, ?)",
		1, settings.Start, settings.Speed, settings.Time, settings.ViewSpeed, settings.TimeCounterDisplay)
	if err != nil {
		log.Fatalf("Error saving settings to database: %v", err)
	}
}

func saveField(id int, fieldText, showField, imgSize string) {
	_, err := db.Exec("REPLACE INTO fields (id, field_text, show_field, img_size) VALUES (?, ?, ?, ?)", id, fieldText, showField, imgSize)
	if err != nil {
		log.Fatalf("Error saving field to database: %v", err)
	}
}

func LoadFields() ([]Field, error) {
	rows, err := db.Query("SELECT id, field_text, show_field, img_size FROM fields")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []Field
	for rows.Next() {
		var field Field
		err := rows.Scan(&field.ID, &field.FieldText, &field.ShowField, &field.ImgSize)
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
