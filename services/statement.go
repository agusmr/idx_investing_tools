package services

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/kevinjanada/idx_investing_tools/models"
)

// StatementService -- Services for creating and modifying financial statements data in DB
type StatementService struct {
	DB *pop.Connection
}

// NewStatementService -- StatementService constructor
func NewStatementService(connectionEnv string) *StatementService {
	db, err := pop.Connect(connectionEnv)
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
	statementRow := s.StatementRowExists(statementRowTitle.ID, stock.ID, date)
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

// StatementRowExists --
func (s *StatementService) StatementRowExists(statementRowTitleID uuid.UUID, stockID uuid.UUID, date time.Time) *models.StatementRow {
	// Check If exists statement row with the same row title, date, and stock id
	statementRow := &models.StatementRow{}
	dateString := date.Format("2006-01-02 15:04:05")
	// dateString := "2017-09-30 00:00:00"
	// fmt.Println(dateString)
	rowTitleIDString := statementRowTitleID.String()
	// rowTitleIDString := "18374f1e-d432-4058-8597-a8ef471eecdf"
	// fmt.Println(rowTitleIDString)
	stockIDString := stockID.String()
	// stockIDString := "b227b28b-3139-458a-b38d-11c64ec1f9f1"
	// fmt.Println(stockIDString)

	queryString := fmt.Sprintf(`
		SELECT  
		sr.id, sr.statement_id, sr.row_description,
		sr.row_order, sr.row_properties, sr.row_title_id
		FROM statement_rows sr
		JOIN statement_row_titles srt
		ON sr.row_title_id = srt.id
		JOIN statement_row_facts srf
		ON srf.statement_row_id = sr.id
		WHERE 
		srt.id = '%s' AND
		srf.date = '%s' AND
		srf.stock_id = '%s';
	`, rowTitleIDString, dateString, stockIDString)

	query := s.DB.RawQuery(queryString)

	err := query.First(statementRow)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(statementRow)

	// fmt.Printf("%v\n", statementRow)
	if statementRow.ID != uuid.Nil {
		fmt.Println("Statement row exists")
		return statementRow
	}
	return nil
}
