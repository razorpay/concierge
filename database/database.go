package database

import (
	"concierge/config"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

//DB ...
var DB *gorm.DB
var err error

//Conn ...
func Conn() {
	log.Info("Creating Connection")
	dbconfig := config.DBConfig

	dbconnURI := dbconfig.DBUsername + ":" + dbconfig.DBPassword + "@tcp(" + dbconfig.Host + ":" + dbconfig.DBPort + ")/" + dbconfig.DBDatabase
	DB, err = gorm.Open("mysql", dbconnURI+"?charset=utf8&parseTime=True&loc=Local")

	DB.DB().SetMaxIdleConns(dbconfig.MaxIdleConns)
	DB.DB().SetMaxOpenConns(dbconfig.MaxOpenConns)
	if err != nil {
		log.Fatal(err)
	}

	DB.LogMode(false)
}

//CloseDB ...
func CloseDB() {
	err := DB.Close()
	if err != nil {
		log.Fatal(err)
	}
}
