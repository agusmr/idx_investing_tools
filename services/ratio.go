package services

import (
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/kevinjanada/idx_investing_tools/tools"
)

type RatioService struct {
	DB *pop.Connection
}

func NewRatioService(connectionEnv string) (*RatioService, error) {
	db, err := pop.Connect(connectionEnv)
	if err != nil {
		return nil, err
	}
	service := &RatioService{DB: db}
	return service, nil
}

func (r *RatioService) CalculateROA(stockCode string, year int) (float64, error) {
	if year == 2017 {
		return 0.0, fmt.Errorf("Could not calculate ROA for 2017, 2016 data not available")
	}
	currentYearDate, err := tools.YearToDate(year)
	if err != nil {
		return 0.0, err
	}
	previousYearDate, err := tools.YearToDate(year - 1)
	if err != nil {
		return 0.0, err
	}

	statementService, err := NewStatementService("tools_development")
	if err != nil {
		return 0.0, err
	}

	currTotalAssets, err := statementService.GetRowFact(
		"Total assets",
		stockCode,
		currentYearDate,
	)
	if err != nil {
		return 0.0, err
	}

	prevTotalAssets, err := statementService.GetRowFact(
		"Total assets",
		stockCode,
		previousYearDate,
	)
	if err != nil {
		return 0.0, err
	}

	currNetIncome, err := statementService.GetRowFact(
		"Total profit (loss)",
		stockCode,
		currentYearDate,
	)
	if err != nil {
		return 0.0, err
	}

	avgTotalAssets := (currTotalAssets.Amount + prevTotalAssets.Amount) / 2
	ROA := currNetIncome.Amount / avgTotalAssets

	return ROA, nil
}
