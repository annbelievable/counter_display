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
	startServer()
}

func startServer() {
	// for static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// for pages
	http.HandleFunc("/", Homepage)

	// for the API
	http.HandleFunc("/latest-counter", LatestCounter)

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

func LastTenCounters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	var c models.CounterLog
	err := row.Scan(&c.Value, &c.Datetime)

	if err != nil {
		return c, err
	}

	return c, nil
}
