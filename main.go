package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var templ = template.Must(template.ParseGlob("form/*"))

type Employee struct{
	Id int
	Name string
	City string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPassword := ""
	dbName :="goblog"
	
	db, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@/"+dbName)

	if err !=nil {
		panic( err.Error());
	}
	return db

}

func Index(w http.ResponseWriter, r *http.Request)  {

	db := dbConn()

	selDb, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
	
	if err !=nil {
		panic(err.Error())
	}

	emp := Employee{}
	res :=[] Employee{}

	for selDb.Next() {

		var id int
		var name, city string
		err := selDb.Scan(&id,&name, &city)

		if err !=nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}
	templ.ExecuteTemplate(w , "Index",res)
	
	defer db.Close();
	
}

func New(w http.ResponseWriter, r *http.Request) {
    templ.ExecuteTemplate(w, "New", nil)
}

func Save(w http.ResponseWriter, r *http.Request){

	db := dbConn()

	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)") 

		if err != nil {
			panic(err.Error());
		}

		insForm.Exec(name,city);
		log.Println("Insert: Name: "+ name +"City: "+ city)
	}

	defer db.Close()
	
	http.Redirect(w, r, "/", 301)

}

func Delete(w http.ResponseWriter, r *http.Request)  {
	db := dbConn()
	prams := mux.Vars(r);
	empId := prams["id"] 
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	
	if err != nil {
		panic(err.Error())
	}

	delForm.Exec(empId)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)

}

func Show(w http.ResponseWriter, r *http.Request){
	db :=dbConn()
	prams := mux.Vars(r)
	empId :=prams["id"]
    selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", empId)
	if err !=nil {
		panic(err.Error())		
	}

	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, city string
		err := selDB.Scan(&id,&name, &city)
		if err !=nil{
			panic(err.Error());
		}
		emp.Id = id
		emp.Name = name
		emp.City = city
	}
	templ.ExecuteTemplate(w,"Show", emp)
	defer db.Close()
}

func Edit(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	prams := mux.Vars(r)
	empId := prams["id"];

	selDB, err :=db.Query("SELECT * FROM Employee where id=?", empId)
	if err !=nil{
		panic(err.Error());
	}
	emp := Employee{}

	for selDB.Next() {
		var id int
		var name, city string
		err := selDB.Scan(&id, &name, &city)
		if err !=nil {
			panic(err.Error())
		} 
		emp.Id = id
		emp.Name = name
		emp.City = city 
	}

	templ.ExecuteTemplate(w,"Edit",emp)
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        city := r.FormValue("city")
        id := r.FormValue("uid")
        insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, city, id)
        log.Println("UPDATE: Name: " + name + " | City: " + city)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
	
	r := mux.NewRouter();
	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/new", New).Methods("GET")
	r.HandleFunc("/new", Save).Methods("POST")
	r.HandleFunc("/delete/{id}", Delete).Methods("GET")
	r.HandleFunc("/show/{id}", Show).Methods("GET")
	r.HandleFunc("/edit/{id}", Edit).Methods("GET")
	r.HandleFunc("/update", Update)
	
	log.Println("Server starting on: http://localhost:8080");
	log.Fatal(http.ListenAndServe(":8080",r))

	
}