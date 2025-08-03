package entity

import "time"

type Account struct {
	ID        uint64 `gorm:"primaryKey"`
	Balance   float64
	CreatedAt time.Time
}

func (Account) TableName() string {
	return "accounts"
}
