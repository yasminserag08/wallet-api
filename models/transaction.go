package models

import "time"

type Transaction struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	WalletID        int       `gorm:"not null" json:"wallet_id"`
	Type            string    `gorm:"not null" json:"type"`
	Amount          int       `gorm:"not null" json:"amount"`
	Category        string    `gorm:"not null" json:"category"`
	Note            string    `gorm:"not null" json:"note"`
	RelatedWalletID int       `json:"related_wallet_id,omitempty"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}
