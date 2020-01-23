package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// Stock is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Stock struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	Code         string       `json:"code" db:"code"`
	Name         string       `json:"name" db:"name"`
	ListingDate  string       `json:"listing_date" db:"listing_date"`
	Shares       int64        `json:"shares" db:"shares"`
	ListingBoard nulls.String `json:"listing_board" db:"listing_board"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (s Stock) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Stocks is not required by pop and may be deleted
type Stocks []Stock

// String is not required by pop and may be deleted
func (s Stocks) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Stock) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Code, Name: "Code"},
		&validators.StringIsPresent{Field: s.Name, Name: "Name"},
		&validators.StringIsPresent{Field: s.ListingDate, Name: "ListingDate"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Stock) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Stock) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
