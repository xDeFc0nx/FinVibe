package types

import (
	"time"
)

type User struct {
	ID                   string                `gorm:"primaryKey"`
	FirstName            string                `json:"FirstName"`
	LastName             string                `json:"lastName"`
	Email                string                `gorm:"type:varchar(100);unique_index"`
	Password             string                `json:"-"`
	CreatedAt            time.Time             `json:"createdAt"`
	UpdatedAt            time.Time             `json:"-"`
	Transactions         []Transaction         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WebSocketConnections []WebSocketConnection `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type WebSocketConnection struct {
	ID           string    `gorm:"primaryKey"`
	ConnectionID string    `gorm:"connectionID"`
	UserID       string    `gorm:"not null"`
	IsActive     bool      `json:"isActive"`
	LastPing     time.Time `json:"lastPing"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Transaction struct {
	ID          string    `gorm:"primaryKey"`
	UserID      string    `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	IsRecurring bool      `gorm:"default:false"`
	Recurring   Recurring `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"-"`
}

type Recurring struct {
	ID            string `gorm:"primaryKey"`
	TransactionID string `gorm:"not null"`

	Frequency string    `gorm:"not null"`
	StartDate time.Time `gorm:"not null"`
	EndDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
