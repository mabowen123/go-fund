package mysql

import "time"

type User struct {
	ID        uint64 `gorm:"primary_key"`
	Username  string `gorm:"type:varchar(16);unique_index"`
	Password  string
	Nickname  string
	CreatedAt time.Time
	UpdatedAt  time.Time
}
