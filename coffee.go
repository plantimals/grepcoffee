package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/plantimals/grepcoffee/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func coffeeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/coffees/")
	log.Print("idStr: " + idStr)
	var c models.Coffee
	id, _ := strconv.Atoi(idStr)
	db.Where(&models.Coffee{ID: uint(id)}).First(&c)
	err := views.ExecuteTemplate(w, "coffee.html", c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// globals
var db, err = gorm.Open("sqlite3", "coffee.db")
var views = template.Must(template.ParseFiles("views/index.html", "views/coffee.html"))
var user *models.User

//var user2 *models.User
var beans *models.Beans
var coffees []*models.Coffee

func main() {
	defer db.Close()
	doMigrations(db)

	var user = models.NewUser("rob", db)
	//var user2 = models.NewUser("not rob", db)
	log.Print(user)
	var beans = models.NewBeans("deathwish", "you'll wish you were dead", db)
	log.Print(beans.Name)

	coffees = append(coffees, models.NewCoffee(user, beans, db))

	http.HandleFunc("/", makeHandler(homeHandler))
	http.HandleFunc("/coffees/", makeHandler(coffeeHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// database functions
func doMigrations(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Coffee{})
	db.AutoMigrate(&models.Beans{})
	db.AutoMigrate(&models.Transition{})
	log.Print("migration done")
}
