package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	"github.com/annbelievable/counter_display/models"
	_ "github.com/mattn/go-sqlite3"
)

var templates = template.Must(template.ParseGlob("./views/*.html"))
var db *sql.DB

func main() {
	fmt.Println("hello")
	db, _ = sql.Open("sqlite3", "./counter.db")
	defer db.Close()
	// go counterLogGenerator()
	startServer()
}

func startServer() {
	// for static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// for pages
	http.HandleFunc("/", Homepage)

	// for the API
	http.HandleFunc("/latest-counter", LatestCounter)
	http.HandleFunc("/last-ten-counter", LastTenCounter)

	http.ListenAndServe(":8080", nil)
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
		http.Error(w, "Something went wrong.", 500)
	}
}

func LatestCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cl, err := selectLatestCounterlog(db)

	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
	}

	data, _ := json.Marshal(&cl)

	w.Write(data)
}

func LastTenCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cls, err := selectLastTenCounterlog(db)

	if err != nil {
		log.Printf("[ERROR] %+v\n", err)
	}

	data, _ := json.Marshal(&cls)

	w.Write(data)
}

func counterLogGenerator() {
	for range time.Tick(time.Second * 5) {
		go func() {
			insertCounterLog(db)
		}()
	}
}

// PRIVATE

func insertCounterLog(db *sql.DB) error {
	value := rand.Intn(9) + 1
	_, err := db.Exec("INSERT INTO counter_log(value, datetime) VALUES (?, ?);", value, time.Now())

	return err
}

func selectLatestCounterlog(db *sql.DB) (models.CounterLog, error) {
	row := db.QueryRow("SELECT value, datetime FROM counter_log ORDER BY datetime DESC LIMIT 0, 1;")
	var cl models.CounterLog
	err := row.Scan(&cl.Value, &cl.Datetime)

	if err != nil {
		return cl, err
	}

	return cl, nil
}

func selectLastTenCounterlog(db *sql.DB) ([]models.CounterLog, error) {
	rows, err := db.Query("SELECT value, datetime FROM counter_log ORDER BY datetime DESC LIMIT 0, 10;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cls []models.CounterLog

	for rows.Next() {
		var cl models.CounterLog
		if err := rows.Scan(&cl.Value, &cl.Datetime); err != nil {
			return cls, err
		}
		cls = append(cls, cl)
	}

	return cls, nil
}
