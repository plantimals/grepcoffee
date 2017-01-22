package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/plantimals/grepcoffee/models"
	"html/template"
	"log"
	"net/http"
)

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func responseWriter(w http.ResponseWriter, view string) {
	err := views.ExecuteTemplate(w, view+".html", coffees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	responseWriter(w, "index")
}

// globals
var db, err = gorm.Open("sqlite3", "coffee.db")
var views = template.Must(template.ParseFiles("views/index.html"))
var user *models.User
var beans *models.Beans
var coffees []*models.Coffee

func main() {
	defer db.Close()
	doMigrations(db)

	var user = models.NewUser("rob", db)
	var beans = models.NewBeans("deathwish", "you'll wish you were dead", db)

	coffees = append(coffees, models.NewCoffee(user, beans, db))
	//coffees[0].Transition(hot, user)
	//coffees = append(coffees, NewCoffee(user, beans))

	http.HandleFunc("/", makeHandler(homeHandler))
	//http.HandleFunc("/coffees/", makeHandler(coffeeHandler))
	http.ListenAndServe(":8080", nil)
}

// database functions
func doMigrations(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Coffee{})
	db.AutoMigrate(&Beans{})
	db.AutoMigrate(&Transition{})
	log.Print("migration done")
}
