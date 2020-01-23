package financialreport

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/pop"
	"github.com/kevinjanada/idx_investing_tools/models"
	"github.com/kevinjanada/idx_investing_tools/tools"
)

type Stock struct {
	Code         string `json:"Code"`
	Name         string `json:"Name"`
	ListingDate  string `json:"ListingDate"`
	Shares       int64  `json:"Shares"`
	ListingBoard string `json:"ListingBoard"`
	Links        []Link `json:"Links"`
}

type Link struct {
	Rel    string `json:"Rel"`
	Href   string `json:"Href"`
	Method string `json:"Method"`
}

type StockAPIResponse struct {
	Draw            int     `json:"draw"`
	RecordsTotal    int     `json:"recordsTotal"`
	RecordsFiltered int     `json:"recordsFiltered"`
	Data            []Stock `json:"data"`
}

// GenerateFetchStockURL --
func GenerateFetchStockURL(start int, length int) string {
	return fmt.Sprintf(
		`https://www.idx.co.id/umbraco/Surface/StockData/GetSecuritiesStock?start=%d&length=%d`,
		start,
		length,
	)
}

// FetchStocksFromDB -- Fetch Stocks Data from local DB
func FetchStocksFromDB() ([]models.Stock, error) {
	tx, err := pop.Connect("api_development")
	if err != nil {
		return nil, err
	}

	stocks := []models.Stock{}
	err = tx.All(&stocks)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

// FetchStocks -- Fetch Stocks Data from IDX API
func FetchStocks() ([]Stock, error) {
	start := 0
	length := 10
	URL := GenerateFetchStockURL(start, length)
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	aggregatedResponse := &StockAPIResponse{}
	tools.JSONToStruct(resp, aggregatedResponse)

	numOfStocksLeft := aggregatedResponse.RecordsTotal
	start += length
	for numOfStocksLeft > 0 {
		URL := GenerateFetchStockURL(start, length)
		resp, err := http.Get(URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		nextResponse := &StockAPIResponse{}
		tools.JSONToStruct(resp, nextResponse)

		aggregatedResponse.Data = append(aggregatedResponse.Data, nextResponse.Data...)

		start += length
		numOfStocksLeft -= length
	}

	return aggregatedResponse.Data, nil
}
