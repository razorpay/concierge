package database

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB
var err error

const (
	//DbUser ...
	DbUser = "mysql"
	//DbPassword ...
	DbPassword = "mysql"
	//DbName ...
	DbName = "concierge"
)

//Conn ...
func Conn() (*gorm.DB, error) {
	db, err = gorm.Open("mysql", "root:helloworld@tcp(192.168.0.102:3306)/concierge?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	return db, nil
}
