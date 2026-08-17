package models

type Wallet struct {
	ID      int `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID  int `gorm:"uniqueIndex,not null" json:"user_id"`
	Balance int `gorm:"not null" json:"balance"`
}
