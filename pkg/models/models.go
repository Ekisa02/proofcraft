package models

import (
	"gorm.io/gorm"
)

// Key stores the public key (private key is not stored by default here)
type Key struct {
	gorm.Model
	Label   string `gorm:"size:255"`
	PubXHex string `gorm:"size:512"`
	PubYHex string `gorm:"size:512"`
}

// Proof stores a recorded proof
type Proof struct {
	gorm.Model
	Label   string `gorm:"size:255"`
	Message string `gorm:"size:1024"`
	RxHex   string `gorm:"size:512"`
	RyHex   string `gorm:"size:512"`
	SHex    string `gorm:"size:512"`
	KeyID   *uint
}
