package domain

import "time"

type Category struct {
	ID        int    `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100);unique"`
	CreatedAt time.Time
}
