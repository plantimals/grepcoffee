package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type State string

const (
	start   State = "start"
	heating State = "heating"
	hot     State = "hot"
	brewing State = "brewing"
	brewed  State = "brewed"
	carafed State = "carafed"
)

type User struct {
	gorm.Model
	Name string `gorm:"not null;unique"`
}

type Transition struct {
	gorm.Model
	From State
	To   State
	Time time.Time
	User *User
}

type Beans struct {
	gorm.Model
	Name string `gorm:"not null;unique"`
	Desc string
}

type Coffee struct {
	gorm.Model
	Name      string
	CurrState State
	Beans     *Beans
	History   []*Transition
}

func NewCoffee(user *User, beans *Beans, db *gorm.DB) *Coffee {
	c := Coffee{Beans: beans, CurrState: start, History: make([]*Transition, 1)}
	db.Create(&c)
	fmt.Println("coffee id: " + string(c.ID))
	//c := new(Coffee)
	//c.Beans = beans
	//c.CurrState = start
	//c.History = make([]*Transition, 1)
	//c.Transition(heating, user)
	return &c
}

func (c *Coffee) Transition(to State, user *User) error {
	tn := new(Transition)
	tn.From = c.CurrState
	tn.To = to
	tn.Time = time.Now()
	tn.User = user
	c.History = append(c.History, tn)
	c.CurrState = to
	c.Name = user.Name + " " + string(c.CurrState) + " " + c.Beans.Name + " @ " + tn.Time.String()
	return nil
}

func NewUser(name string, db *gorm.DB) *User {
	var u *User
	err := db.Where(&User{Name: name}).First(u)
	if err != nil {
		fmt.Println("caught the missing where")
	}
	//err = db.Create(&User{Name: name})
	//if err != nil {
	//    fmt.Println("caught the uniqueness")
	//}
	fmt.Println("====")
	return u
}

func NewBeans(name string, desc string, db *gorm.DB) *Beans {
	b := new(Beans)
	b.Name = name
	b.Desc = desc
	return b
}
