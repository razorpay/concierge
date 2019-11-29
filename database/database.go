package database

import (
	"concierge/config"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB
var err error

//Conn ...
func Conn() (*gorm.DB, error) {

	dbconfig := config.DBConfig

	dbconnURI := dbconfig.DBUsername + ":" + dbconfig.DBPassword + "@tcp(" + dbconfig.Host + ":" + dbconfig.DBPort + ")"
	db, err = gorm.Open("mysql", dbconnURI+"/concierge?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	return db, nil
}
