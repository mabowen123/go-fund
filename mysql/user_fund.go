package mysql

import (
	"time"
)

type UserFund struct {
	ID        uint `gorm:"primary_key"`
	UserId    int
	FundId    string
	Share     float64
	Amount    float64
	CreatedAt time.Time
	UpdatedAt  time.Time
}

