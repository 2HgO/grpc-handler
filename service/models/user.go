package models

import (
	"github.com/go-bongo/bongo"
)

// User ...
type User struct {
	bongo.DocumentBase `json:",inline" bson:",inline" structs:"-"`
	FirstName          string `json:"firstName,omitempty" bson:"firstName" structs:"firstName,omitempty"`
	LastName           string `json:"lastName,omitempty" bson:"lastName" structs:"lastName,omitempty"`
}
