package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"fmt"
	"html/template"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
	err error
)

type Cost struct {
	Id int64 `json:"id"`
	ElectricAmount int64 `json:"electric_amount"`
	ElectricPrice float64 `json:"electric_price"`
	WaterAmount int64 `json:"water_amount"`
	WaterPrice float64 `json:"water_price"`
	CheckedDate string `json:"checked_date"`
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
	// route
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(":3333", nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	rows, err := db.Query("SELECT * FROM cost")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var funcMap = template.FuncMap{
		"multiplication": func(n float64, f float64) float64 {
			return n * f
		},
		"addOne": func(n int) int {
			return n + 1
		},
	}
	var costs []Cost
	var cost Cost
	for rows.Next() {
		err = rows.Scan(&cost.Id, &cost.ElectricAmount,
			&cost.ElectricPrice, &cost.WaterAmount, &cost.WaterPrice, &cost.CheckedDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		costs = append(costs, cost)
	}
	//t, err := template.ParseFiles("tmpl/list.html")
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = t.Execute(w, costs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "tmpl/create.html")
		return
	}
	var cost Cost
	cost.ElectricAmount, _ = strconv.ParseInt(r.FormValue("ElectricAmount"), 10, 64)
	cost.ElectricPrice, _ = strconv.ParseFloat(r.FormValue("ElectricPrice"), 64)
	cost.WaterAmount, _ = strconv.ParseInt(r.FormValue("WaterAmount"), 10, 64)
	cost.WaterPrice, _ = strconv.ParseFloat(r.FormValue("WaterPrice"), 64)
	cost.CheckedDate = r.FormValue("CheckedDate")
	fmt.Println(cost)

	// Save to database
	stmt, err := db.Prepare(`
		INSERT INTO cost(electric_amount, electric_price, water_amount, water_price, checked_date)
		VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Println("Prepare query error")
		panic(err)
	}
	_, err = stmt.Exec(cost.ElectricAmount, cost.ElectricPrice,
				cost.WaterAmount, cost.WaterPrice, cost.CheckedDate)
	if err != nil {
		fmt.Println("Execute query error")
		panic(err)
	}
	http.Redirect(w, r, "/list", 301)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not allowed", http.StatusBadRequest)
	}
	var cost Cost
	cost.Id, _ = strconv.ParseInt(r.FormValue("Id"), 10, 64)
	cost.ElectricAmount, _ = strconv.ParseInt(r.FormValue("ElectricAmount"), 10, 64)
	cost.ElectricPrice, _ = strconv.ParseFloat(r.FormValue("ElectricPrice"), 64)
	cost.WaterAmount, _ = strconv.ParseInt(r.FormValue("WaterAmount"), 10, 64)
	cost.WaterPrice, _ = strconv.ParseFloat(r.FormValue("WaterPrice"), 64)
	cost.CheckedDate = r.FormValue("CheckedDate")
	fmt.Println(cost)
	stmt, err := db.Prepare(`
		UPDATE cost SET electric_amount=?, electric_price=?, water_amount=?, water_price=?, checked_date=?
		WHERE id=?
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res, err := stmt.Exec(cost.ElectricAmount, cost.ElectricPrice,
		cost.WaterAmount, cost.WaterPrice, cost.CheckedDate, cost.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/list", 301)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	var costId, _ = strconv.ParseInt(r.FormValue("Id"), 10, 64)
	stmt, err := db.Prepare("DELETE FROM cost WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res, err := stmt.Exec(costId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/list", 301)

}