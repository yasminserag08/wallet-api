package models

type Wallet struct {
	ID      uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID  uint `gorm:"uniqueIndex,not null" json:"user_id"`
	Balance int  `gorm:"not null" json:"balance"`
}
