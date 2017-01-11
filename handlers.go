package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strconv"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "tmpl/register.html")
		return
	}
	// grab user info
	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")
	// Check existence of user
	var user User
	err := db.QueryRow("SELECT username, password, role FROM users WHERE username=?",
		username).Scan(&user.Username, &user.Password, &user.Role)
	switch {
	// user is available
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkInternalServerError(err, w)
		// insert to database
		_, err = db.Exec(`INSERT INTO users(username, password, role) VALUES(?, ?, ?)`,
			username, hashedPassword, role)
		fmt.Println("Created user: ", username)
		checkInternalServerError(err, w)
	case err != nil:
		http.Error(w, "loi: "+err.Error(), http.StatusBadRequest)
		return
	default:
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "tmpl/login.html")
		return
	}
	// grab user info from the submitted form
	username := r.FormValue("usrname")
	password := r.FormValue("psw")
	// query database to get match username
	var user User
	err := db.QueryRow("SELECT username, password FROM users WHERE username=?",
		username).Scan(&user.Username, &user.Password)
	checkInternalServerError(err, w)
	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/login", 301)
	}
	authenticated = true
	http.Redirect(w, r, "/list", 301)

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	authenticated = false
	isAuthenticated(w, r)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated(w, r)
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
	isAuthenticated(w, r)
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
	isAuthenticated(w, r)
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
	isAuthenticated(w, r)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated(w, r)
	http.Redirect(w, r, "/list", 301)
}