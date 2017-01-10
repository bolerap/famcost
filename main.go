package main

import (
	"database/sql"
	"net/http"
	"os"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db  *sql.DB
	err error
)

type Cost struct {
	Id             int64   `json:"id"`
	ElectricAmount int64   `json:"electric_amount"`
	ElectricPrice  float64 `json:"electric_price"`
	WaterAmount    int64   `json:"water_amount"`
	WaterPrice     float64 `json:"water_price"`
	CheckedDate    string  `json:"checked_date"`
}

func main() {
	db, err = sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// test connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	os.Setenv("PORT", "8898")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	// route
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.Handle("/statics/",
		http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))),
	)
	http.ListenAndServe(":" + port, nil)
}
