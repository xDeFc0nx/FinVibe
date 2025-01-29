package types

import (
	"time"
)

type WebSocketConnection struct {
	ID           string    `gorm:"primaryKey"`
	ConnectionID string    `gorm:"connectionID"`
	UserID       string    `gorm:"not null;index"`
	IsActive     bool      `                      json:"isActive"`
	LastPing     time.Time `                      json:"lastPing"`
	CreatedAt    time.Time `                      json:"createdAt"`
}

type User struct {
	ID                   string                `gorm:"primaryKey;index"`
	FirstName            string                `                                                                       json:"FirstName"`
	LastName             string                `                                                                       json:"lastName"`
	Email                string                `gorm:"type:varchar(100);unique_index"`
	Password             string                `                                                                       json:"Password"`
	Currency             string                `                                                                       json:"Currency"`
	CreatedAt            time.Time             `                                                                       json:"createdAt"`
	UpdatedAt            time.Time             `                                                                       json:"-"`
	Accounts             []Accounts            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	WebSocketConnections []WebSocketConnection `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Accounts struct {
	ID           string        `gorm:"primaryKey;column:id"`
	UserID       string        `gorm:"not null;column:user_id"`
	User         User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type         string        `                                                                          json:"type"`
	Income       float64       `                                                                          json:"income"`
	Expense      float64       `                                                                          json:"expense"`
	Balance      float64       `                                                                          json:"balance"`
	Transactions []Transaction `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time     `                                                                          json:"createdAt"`
	UpdatedAt    time.Time     `                                                                          json:"-"`
}

type Transaction struct {
	ID        string   `gorm:"primaryKey;column:id"`
	UserID    string   `gorm:"not null;column:user_id"`
	User      User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AccountID string   `gorm:"not null;column:account_id"`
	Account   Accounts `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type      string   `                                                                          json:"type"`
	Amount    float64  `                                                                          json:"amount"`

	Description string    `json:"description"`
	IsRecurring bool      `                   gorm:"default:false"`
	Recurring   Recurring `                   gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"-"`
}

type Recurring struct {
	ID            string     `gorm:"primaryKey"`
	TransactionID string     `gorm:"not null"`
	Amount        float64    `                  json:"amount"`
	Frequency     string     `                  json:"frequency"`
	StartDate     time.Time  `                  json:"startDate"`
	NextDate      time.Time  `                  json:"nextDate"`
	EndDate       *time.Time `                  json:"endDate"`
	CreatedAt     time.Time  `                  json:"createdAt"`
	UpdatedAt     time.Time  `                  json:"updatedAt"`
}

type Budget struct {
	ID          string    `gorm:"primaryKey;column:id"`
	UserID      string    `gorm:"not null;column:user_id"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AccountID   string    `gorm:"not null;column:account_id"`
	Account     Accounts  `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TotalSpent  float64   `gorm:"column:total_spent"                                                 json:"totalSpent"`
	Limit       float64   `gorm:"column:limit"                                                       json:"limit"`
	Description string    `                                                                          json:"description"`
	CreatedAt   time.Time `                                                                          json:"createdAt"`
	UpdatedAt   time.Time `                                                                          json:"-"`
}

type Goal struct {
	ID          string    `gorm:"primaryKey;column:id"`
	UserID      string    `gorm:"not null;column:user_id"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AccountID   string    `gorm:"not null;column:account_id"`
	Account     Accounts  `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GoalAmount  float64   `gorm:"goal_amount"`
	Amount      float64   `gorm:"amount"`
	Description string    `                                                                          json:"description"`
	CreatedAt   time.Time `                                                                          json:"createdAt"`
	UpdatedAt   time.Time `                                                                          json:"-"`
}
