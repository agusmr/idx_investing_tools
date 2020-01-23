package services

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/kevinjanada/idx_investing_tools/models"
)

// StatementFinancialPosition -- Constant
const StatementFinancialPosition = "financial position"

// StatementProfitOrLoss -- Constant
const StatementProfitOrLoss = "profit or loss"

// RowTotalAssets -- Constant
const RowTotalAssets = "total assets"

// RowNetIncome -- Constant
const RowNetIncome = "net income"

// StatementService -- Services for creating and modifying financial statements data in DB
type StatementService struct {
	DB *pop.Connection
}

// NewStatementService -- StatementService constructor
func NewStatementService() *StatementService {
	db, err := pop.Connect("tools_development")
	if err != nil {
		fmt.Println(err)
	}
	service := &StatementService{DB: db}
	return service
}

// NewStatement -- Create a new statement name in DB
func (s *StatementService) NewStatement(statementName string) error {
	statement := &models.Statement{}
	s.DB.Where("name = $1", statementName).First(statement)

	if statement.ID != uuid.Nil {
		return fmt.Errorf("Statement with that name already exists")
	}

	statement.Name = statementName
	err := s.DB.Save(statement)
	if err != nil {
		return err
	}
	return nil
}

// NewStatementRowTitle -- Create a new statement row title in DB
func (s *StatementService) NewStatementRowTitle(rowTitle string) error {
	statementRowTitle := &models.StatementRowTitle{}
	s.DB.Where("name = $1", rowTitle).First(statementRowTitle)

	if statementRowTitle.ID != uuid.Nil {
		return fmt.Errorf("Row title with that name already exists")
	}

	statementRowTitle.Title = rowTitle
	err := s.DB.Save(statementRowTitle)
	if err != nil {
		return err
	}
	return nil
}

// InsertUpdateStatementRow -- Insert information row to a stock's statement
func (s *StatementService) InsertUpdateStatementRow(
	stockCode string,
	statementName string,
	rowTitle string,
	rowAmount float64,
	date time.Time,
) error {
	stock := &models.Stock{}
	s.DB.Where("code = $1", stockCode).First(stock)
	if stock.ID == uuid.Nil {
		return fmt.Errorf("Stock not found")
	}

	statement := &models.Statement{}
	s.DB.Where("name = $1", statementName).First(statement)
	if statement.ID == uuid.Nil {
		return fmt.Errorf("Statement name not found")
	}

	statementRowTitle := &models.StatementRowTitle{}
	s.DB.Where("title = $1", rowTitle).First(statementRowTitle)
	if statementRowTitle.ID == uuid.Nil {
		return fmt.Errorf("Row title not found")
	}

	// Check If exists statement row with the same row title, date, and stock id
	statementRow := s.statementRowExists(statementRowTitle.ID, stock.ID, date)
	// If exists, Update the statementRowFact
	if statementRow != nil {
		statementRowFact := &models.StatementRowFact{}
		s.DB.Where("statement_row_id = $1", statementRow.ID)

		statementRowFact.Amount = rowAmount
		s.DB.Save(statementRowFact)
		return nil
	}

	// If statementRow does not exists, create a new one
	statementRow = &models.StatementRow{
		StatementID: statement.ID,
		RowTitleID:  statementRowTitle.ID,
	}

	err := s.DB.Save(statementRow)
	if err != nil {
		return err
	}

	err = s.DB.Where("statement_id = $1", statement.ID).
		Where("row_title_id = $1", statementRowTitle.ID).
		First(statementRow)

	statementRowFact := &models.StatementRowFact{
		StatementRowID: statementRow.ID,
		StockID:        stock.ID,
		Date:           date,
		Amount:         rowAmount,
	}

	err = s.DB.Save(statementRowFact)
	if err != nil {
		return err
	}
	return nil
}

func (s *StatementService) statementRowExists(statementRowTitleID uuid.UUID, stockID uuid.UUID, date time.Time) *models.StatementRow {
	// Check If exists statement row with the same row title, date, and stock id
	statementRow := &models.StatementRow{}
	s.DB.RawQuery(`
		SELECT  
		sr.id, sr.statement_id, sr.row_description,
		sr.row_order, sr.row_properties, sr.row_title_id
		FROM statement_rows sr
		JOIN statement_row_titles srt
		ON sr.row_title_id = srt.id
		JOIN statement_row_facts srf
		ON srf.statement_row_id = sr.id
		WHERE 
		srt.id = $1 AND
		srf.date = $2 AND
		srf.stock_id = $3
	`, statementRowTitleID, date, stockID)

	if statementRow.ID != uuid.Nil {
		return statementRow
	}
	return nil
}
