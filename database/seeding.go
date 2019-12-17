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
		Name:     "ankit.infra",
		Username: "ankit.infra",
		Email:    "ankit.infra@razorpay.com",
		Admin:    1,
	})
}
