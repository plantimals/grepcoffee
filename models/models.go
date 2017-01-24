package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

//type State string

const (
	Start   = "start"
	Heating = "heating"
	Hot     = "hot"
	Brewing = "brewing"
	Brewed  = "brewed"
	Carafed = "carafed"
)

type User struct {
	gorm.Model
	Name string `gorm:"not null;unique"`
}

type Transition struct {
	gorm.Model
	From string
	To   string
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
	ID        uint `gorm:"primary_key"`
	Name      string
	CurrState string
	Beans     *Beans
	History   []*Transition
}

func NewCoffee(user *User, beans *Beans, db *gorm.DB) *Coffee {
	var c Coffee
	db.Create(&Coffee{Name: "", Beans: beans, CurrState: Start})
	db.Where(&Coffee{CurrState: Start}).First(&c)
	c.History = make([]*Transition, 100)
	c.Transition(Heating, user, db)
	return &c
}

func (c *Coffee) Transition(to string, user *User, db *gorm.DB) error {
	db.Create(&Transition{From: c.CurrState, To: to, Time: time.Now(), User: user})
	var tn Transition
	db.Where(&Transition{From: c.CurrState, To: to, User: user}).First(&tn)
	tn.From = c.CurrState
	tn.To = to
	tn.Time = time.Now()
	tn.User = user
	c.History = append(c.History, &tn)
	c.CurrState = to
	c.Name = c.MkName(user)
	return nil
}

func (c *Coffee) MkName(user *User) string {
	t := time.Now().Local()
	return fmt.Sprintf("%s @ %02d:%02d %s", user.Name, t.Hour(), t.Minute(), t.Weekday())
}

func (u *User) String() string {
	return u.Name
}

func NewUser(name string, db *gorm.DB) *User {
	var u User
	db.Where(&User{Name: name}).First(&u)
	if u.Name == name {
		return &u
	}
	db.Create(&User{Name: name})
	db.Where(&User{Name: name}).First(&u)
	return &u
}

func NewBeans(name string, desc string, db *gorm.DB) *Beans {
	var b Beans
	db.Where(&Beans{Name: name}).First(&b)
	if b.Name == name {
		log.Print("found beans")
		return &b
	}
	db.Create(&Beans{Name: name, Desc: desc})
	db.Where(&Beans{Name: name}).First(&b)
	return &b
}
