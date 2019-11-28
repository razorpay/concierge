package models

import (
	"time"
)

//Users ...
type Users struct {
	ID              uint `gorm:"type:bigint(20) unsigned auto_increment;primary_key"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `gorm:"unique_index:users_email_username_deleted_at_unique"`
	Name            string     `gorm:"column:name;type:varchar(255);not null" form:"name" binding:"required,min=3,max=100"`
	Username        string     `gorm:"column:username;type:varchar(255);not null;unique_index:users_email_username_deleted_at_unique,users_username_unique"`
	Email           string     `gorm:"column:email;type:varchar(255);not null;unique_index:users_email_username_deleted_at_unique" form:"email" binding:"required,email"`
	EmailVerifiedAt time.Time  `gorm:"column:email_verified_at;type:timestamp;default:null"`
	Admin           int        `gorm:"column:admin;type:tinyint(4);default:null" form:"admin"`
	GoogleID        string     `gorm:"column:google_id;type:varchar(255);not null;default:''"`
	AccessToken     string     `gorm:"column:access_token;type:varchar(255);not null;default:''"`
	Password        string     `gorm:"column:password;type:varchar(255);default:null"`
	RememberToken   string     `gorm:"column:remember_token;type:varchar(255);default:null"`
}
