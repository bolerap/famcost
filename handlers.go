package main

import (
	"net/http"
	"fmt"
	"html/template"
	"strconv"
)

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
	rows, err := db.Query("SELECT * FROM cost")
	checkInternalServerError(err, w)
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
		checkInternalServerError(err, w)
		costs = append(costs, cost)
	}
	t, err := template.New("list.html").Funcs(funcMap).ParseFiles("tmpl/list.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, costs)
	checkInternalServerError(err, w)

}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
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
	http.Redirect(w, r, "/", 301)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
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
	checkInternalServerError(err, w)
	res, err := stmt.Exec(cost.ElectricAmount, cost.ElectricPrice,
		cost.WaterAmount, cost.WaterPrice, cost.CheckedDate, cost.Id)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", 301)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	var costId, _ = strconv.ParseInt(r.FormValue("Id"), 10, 64)
	stmt, err := db.Prepare("DELETE FROM cost WHERE id=?")
	checkInternalServerError(err, w)
	res, err := stmt.Exec(costId)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)
	http.Redirect(w, r, "/", 301)

}
