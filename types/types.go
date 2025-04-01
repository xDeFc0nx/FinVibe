package types

import (
	"time"
)

type WebSocket struct {
	ID           string    `db:"id" json:"id"`
	ConnectionID string    `db:"connection_id" json:"connectionId"`
	UserID       string    `db:"user_id" json:"userId"`
	IsActive     bool      `db:"is_active" json:"isActive"`
	LastPing     time.Time `db:"last_ping" json:"lastPing"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}

type User struct {
	ID        string    `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"firstName"`
	LastName  string    `db:"last_name" json:"lastName"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	Currency  string    `db:"currency" json:"currency"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type Accounts struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"userId"`
	Type      string    `db:"type" json:"type"`
	Income    float64   `db:"income" json:"income"`
	Expense   float64   `db:"expense" json:"expense"`
	Balance   float64   `db:"balance" json:"balance"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type Transaction struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"userId"`
	AccountID   string    `db:"account_id" json:"accountId"`
	Type        string    `db:"type" json:"type"`
	Amount      float64   `db:"amount" json:"amount"`
	Description string    `db:"description" json:"description"`
	IsRecurring bool      `db:"is_recurring" json:"isRecurring"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

type Recurring struct {
	ID            string     `db:"id" json:"id"`
	TransactionID string     `db:"transaction_id" json:"transactionId"`
	Amount        float64    `db:"amount" json:"amount"`
	Frequency     string     `db:"frequency" json:"frequency"`
	StartDate     time.Time  `db:"start_date" json:"startDate"`
	NextDate      time.Time  `db:"next_date" json:"nextDate"`
	EndDate       *time.Time `db:"end_date" json:"endDate,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updatedAt"`
}

type Budget struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"userId"`
	AccountID   string    `db:"account_id" json:"accountId"`
	TotalSpent  float64   `db:"total_spent" json:"totalSpent"`
	LimitAmount float64   `db:"limit_amount" json:"limit"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

type Goal struct {
	ID            string  `db:"id" json:"id"`
	UserID        string  `db:"user_id" json:"userId"`
	AccountID     string  `db:"account_id" json:"accountId"`
	GoalAmount    float64 `db:"goal_amount" json:"goalAmount"`
	CurrentAmount float64 `db:"current_amount" json:"currentAmount"`
	Description   string  `db:"description" json:"description"`
}
