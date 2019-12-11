package database

import (
	"concierge/models"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

//Seeding ...
func Seeding() {
	var db *gorm.DB
	db, err := Conn()
	if err != nil {
		log.Error(err)
	}
	db.FirstOrCreate(&models.Users{}, &models.Users{
		Name:     "ankit.infra",
		Username: "ankit.infra",
		Email:    "ankit.infra@razorpay.com",
		Admin:    1,
	})
	defer db.Close()
}
