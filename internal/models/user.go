package models

import (
	"time"
)

type User struct {
	Id        int64     `json:"id,omitempty" db:"id"`
	Name      string    `json:"name,omitempty" db:"name"`
	Email     string    `json:"email,omitempty" db:"email"`
	Password  string    `json:"password,omitempty" db:"password"`
	PhoneNo   string    `json:"phoneNo,omitempty" db:"phone_no"`
	RoleId    int64     `json:"roleId,omitempty" db:"role_id"`
	CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at"`
}

type ctxKey string

const (
	UserCtxKey ctxKey = ctxKey("user")
)
