package models

import (
	"github.com/go-bongo/bongo"
)

// User ...
type User struct {
	bongo.DocumentBase `json:",inline" structs:"-"`
	FirstName          string `json:"firstName,omitempty" binding:"required"`
	LastName           string `json:"lastName,omitempty" binding:"required"`
}
