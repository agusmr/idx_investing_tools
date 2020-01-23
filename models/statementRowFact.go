package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"time"
	"github.com/gobuffalo/validate/validators"
)
// StatementRowFact is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type StatementRowFact struct {
    ID uuid.UUID `json:"id" db:"id"`
    StatementRowID uuid.UUID `json:"statement_row_id" db:"statement_row_id"`
    StockID uuid.UUID `json:"stock_id" db:"stock_id"`
    Date time.Time `json:"date" db:"date"`
    Amount float64 `json:"amount" db:"amount"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (s StatementRowFact) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// StatementRowFacts is not required by pop and may be deleted
type StatementRowFacts []StatementRowFact

// String is not required by pop and may be deleted
func (s StatementRowFacts) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *StatementRowFact) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: s.Date, Name: "Date"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *StatementRowFact) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *StatementRowFact) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
