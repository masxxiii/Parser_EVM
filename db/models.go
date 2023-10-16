package db

import "time"

// Block struct provides model definition for the block table.
type Block struct {
	ID        uint   `gorm:"primaryKey"`
	Network   string `gorm:"unique"`
	Height    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Wallet struct provides model definition for the wallet table.
type Wallet struct {
	ID        uint   `gorm:"primaryKey"`
	Address   string `gorm:"unique"`
	CreatedAt time.Time
}

// Transaction struct provides model definition for the transaction table.
type Transaction struct {
	ID        uint   `gorm:"primaryKey"`
	Hash      string `gorm:"unique"`
	From      string
	To        string
	Value     string
	CreatedAt time.Time
}
