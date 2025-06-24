package domain

import "time"

type StockMovement struct {
	ID        int `gorm:"primaryKey"`
	ProductID int
	UserID    int
	Type      string `gorm:"type:enum('in','out')"`
	Quantity  int
	Note      string
	CreatedAt time.Time

	Product Product `gorm:"foreignKey:ProductID"`
	User    User    `gorm:"foreignKey:UserID"`
}
