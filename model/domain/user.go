package domain

import "time"

type User struct {
	ID        int    `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100)"`
	Email     string `gorm:"type:varchar(100);unique"`
	Password  string `gorm:"type:text"`
	Role      string `gorm:"type:enum('admin','staff')"`
	CreatedAt time.Time
}
