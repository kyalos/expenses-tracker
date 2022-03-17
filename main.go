package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Expenses struct {
	Expense_id    int
	Expense_name  string
	Expense_value string
	Incurred_on   string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "kyalos"
	dbPass := "sealed"
	dbName := "expenses_tracker"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Expenses ORDER BY expense_id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Expenses{}
	res := []Expenses{}

	for selDB.Next() {
		var expense_id int
		var expense_name, expense_value, incurred_on string
		err = selDB.Scan(&expense_id, &expense_name, &expense_value, &incurred_on)
		if err != nil {
			panic(err.Error())
		}
		emp.Expense_id = expense_id
		emp.Expense_name = expense_name
		emp.Expense_value = expense_value
		emp.Incurred_on = incurred_on
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("expense_id")
	selDB, err := db.Query("SELECT * FROM Expenses WHERE expense_id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Expenses{}
	for selDB.Next() {
		var expense_id int
		var expense_name, expense_value, incurred_on string
		err = selDB.Scan(&expense_id, &expense_name, &expense_value, &incurred_on)
		if err != nil {
			panic(err.Error())
		}
		emp.Expense_id = expense_id
		emp.Expense_name = expense_name
		emp.Expense_value = expense_value
		emp.Incurred_on = incurred_on
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("expense_id")
	selDB, err := db.Query("SELECT * FROM Expenses WHERE expense_id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Expenses{}
	for selDB.Next() {
		var expense_id int
		var expense_name, expense_value, incurred_on string
		err = selDB.Scan(&expense_id, &expense_name, &expense_value, &incurred_on)
		if err != nil {
			panic(err.Error())
		}
		emp.Expense_id = expense_id
		emp.Expense_name = expense_name
		emp.Expense_value = expense_value
		emp.Incurred_on = incurred_on
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		expense_name := r.FormValue("expense_name")
		expense_value := r.FormValue("expense_value")
		incurred_on := r.FormValue("incurred_on")
		insForm, err := db.Prepare("INSERT INTO Expenses(expense_name, expense_value, incurred_on) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(expense_name, expense_value, incurred_on)
		log.Println("INSERT: Expense_name: " + expense_name + " | Expense_value: " + expense_value + " | Incurred_on: " + incurred_on)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		expense_name := r.FormValue("expense_name")
		expense_value := r.FormValue("expense_value")
		incurred_on := r.FormValue("incurred_on")
		expense_id := r.FormValue("expense_id")

		//expense_id := r.URL.Query().Get("expense_id")
		insForm, err := db.Prepare("UPDATE Expenses SET expense_name=?, expense_value=?, incurred_on=? WHERE expense_id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(expense_name, expense_value, expense_id, incurred_on)
		log.Println("UPDATE: Expense_name: " + expense_name + " | Expense_value: " + expense_value + " | Incurred_on: " + incurred_on + " | Where: " + expense_id)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("expense_id")
	delForm, err := db.Prepare("DELETE FROM Expenses WHERE expense_id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8081")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8081", nil)
}
