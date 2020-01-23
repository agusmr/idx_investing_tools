package models

import (
	"encoding/json"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"time"
)
// StatementRow is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type StatementRow struct {
    ID uuid.UUID `json:"id" db:"id"`
    StatementID uuid.UUID `json:"statement_id" db:"statement_id"`
    RowTitleID uuid.UUID `json:"row_title_id" db:"row_title_id"`
    RowOrder nulls.Int `json:"row_order" db:"row_order"`
    RowDescription nulls.Int `json:"row_description" db:"row_description"`
    RowProperties nulls.Int `json:"row_properties" db:"row_properties"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (s StatementRow) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// StatementRows is not required by pop and may be deleted
type StatementRows []StatementRow

// String is not required by pop and may be deleted
func (s StatementRows) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *StatementRow) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *StatementRow) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *StatementRow) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
