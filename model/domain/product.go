package domain

import "time"

type Product struct {
	ID         int    `gorm:"primaryKey"`
	Name       string `gorm:"type:varchar(100)"`
	CategoryID int
	Stock      int
	CreatedAt  time.Time

	Category Category `gorm:"foreignKey:CategoryID"`
}
