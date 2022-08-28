package database

import (
	"concierge/models"
)

//Seeding ...
func Seeding() {
	if DB == nil {
		Conn()
	}
	DB.FirstOrCreate(&models.Users{}, &models.Users{
		Name:     "admin",
		Username: "admin",
		Email:    "admin@razorpay.com",
		Admin:    1,
	})
}
