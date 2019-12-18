package models

import (
	"time"
)

//Leases ...
type Leases struct {
	ID          uint `gorm:"type:int(20) unsigned auto_increment;primary_key;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	UserID      uint   `gorm:"column:user_id;type:int(20) unsigned;not null;"`
	User        Users  `gorm:"foreignkey:UserID"`
	GroupID     string `gorm:"column:group_id;type:varchar(20);default:null"`
	LeaseIP     string `gorm:"column:lease_ip;type:varchar(255);not null"`
	LeaseType   string `gorm:"column:lease_type;type:varchar(255);not null"`
	Protocol    string `gorm:"column:protocol;type:varchar(255);default:null"`
	PortFrom    string `gorm:"column:port_from;type:varchar(255);default:null"`
	PortTo      string `gorm:"column:port_to;type:varchar(255);default:null"`
	Expiry      uint   `gorm:"column:expiry;type:int(10) unsigned;not null"`
	InviteEmail string `gorm:"column:invite_email;type:varchar(255);default:null"`
}
