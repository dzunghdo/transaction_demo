package entity

import "time"

type Transaction struct {
	ID                   uint64 `gorm:"primaryKey;autoIncrement"`
	SourceAccountID      uint64
	DestinationAccountID uint64
	Amount               float64
	TransactionTime      time.Time

	// Relationships
	SourceAccount      Account `gorm:"foreignKey:SourceAccountID"`
	DestinationAccount Account `gorm:"foreignKey:DestinationAccountID"`
}

func (Transaction) TableName() string {
	return "transactions"
}
