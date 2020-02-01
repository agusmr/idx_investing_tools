package handlers

import (
	"fmt"
	"github.com/kevinjanada/idx_investing_tools/services"
	"github.com/kevinjanada/idx_investing_tools/tools"
)

func CalculateRatios(year int) error {
	ratioService, err := services.NewRatioService("tools_development")
	if err != nil {
		return err
	}
	stocksService, err := services.NewStockService("tools_development")
	if err != nil {
		return err
	}
	statementService, err := services.NewStatementService("tools_development")
	if err != nil {
		return err
	}

	stocks, err := stocksService.FetchStocksFromDB()
	if err != nil {
		return err
	}
	for _, s := range stocks {
		if s.ListingDate == "-" {
			continue
		}
		// Calculate ROA ----------
		ROAAmount, err := ratioService.CalculateROA(s.Code, year)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				continue
			}
			return err
		}
		date, err := tools.YearToDate(year)
		if err != nil {
			return err
		}
		statementService.InsertUpdateStatementRow(
			s.Code,
			"ratios",
			"ROA",
			ROAAmount,
			date,
		)
		// TODO: Calculate other ratios here
		// ------
		// ------
	}

	fmt.Println("Finished calculating ROA")
	return nil
}
