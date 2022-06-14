package models

import (
	"github.com/revel/revel"
)

// User represents user data
// swagger:model User
type User struct {
	// swagger:ignore
	ID int64 `json:"id"`

	// Username
	// required: true
	// min length: 4
	// max length: 15
	// example: mkulish
	Username string  `json:"username"`

	// Password
	// required: true
	// min length: 5
	// max length: 15
	// example: test123!
	Password string  `json:"password"`

	// swagger:ignore
	HashedPassword []byte `json:"-"`
}

// LoginResponse represents login response JSON
// swagger:model LoginResponse
type LoginResponse struct {
	// Operation success flag
	// required: true
	// type: boolean
	OK bool `json:"ok"`

	// Auth token
	// required: true
	// type: string
	Token string `json:"token"`
}

// Validate checks user data
func (u *User) Validate(v *revel.Validation) {
	v.Check(u.Username,
		revel.Required{},
		revel.MinSize{Min: 4},
		revel.MaxSize{Max: 15},
	).Key("username")

	v.Check(u.Password,
		revel.Required{},
		revel.MinSize{Min: 5},
		revel.MaxSize{Max: 15},
	).Key("password")
}
