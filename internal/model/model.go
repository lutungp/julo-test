package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid"`
	Title    string
	SubTitle string
	Text     string
}

type Account struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid"`
	CustomerXid string
	Status      string
	DisabledAt  time.Time
}

type Wallet struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid"`
	CustomerId string
	Status     string
	EnableAt   time.Time
	Balance    float64
	DisabledAt time.Time
}

type FillWallet struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid"`
	CustomerId  string
	WalletId    string
	Value       float64
	ReferenceId string
}
