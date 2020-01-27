package handlers

import (
	//"fmt"
	"github.com/kevinjanada/idx_investing_tools/services"
)

func UpdateStockList() error {
	stockService, err := services.NewStockService("tools_development")
	if err != nil {
		return err
	}
	stocksData, err := stockService.FetchStocks()
	if err != nil {
		return err
	}
	err = stockService.SaveStockDataToDB(stocksData)
	if err != nil {
		return err
	}
	return nil
}
