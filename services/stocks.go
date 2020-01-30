package services

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/kevinjanada/idx_investing_tools/models"
	"github.com/kevinjanada/idx_investing_tools/tools"
)

type StockService struct {
	DB *pop.Connection
}

func NewStockService(connectionKey string) (*StockService, error) {
	db, err := pop.Connect(connectionKey)
	if err != nil {
		return nil, err
	}
	service := &StockService{DB: db}
	return service, nil
}

type StockData struct {
	Code         string       `json:"Code"`
	Name         string       `json:"Name"`
	ListingDate  string       `json:"ListingDate"`
	Shares       int64        `json:"Shares"`
	ListingBoard nulls.String `json:"ListingBoard"`
	Links        []Link       `json:"Links"`
}

type Link struct {
	Rel    string `json:"Rel"`
	Href   string `json:"Href"`
	Method string `json:"Method"`
}

type StockAPIResponse struct {
	Draw            int         `json:"draw"`
	RecordsTotal    int         `json:"recordsTotal"`
	RecordsFiltered int         `json:"recordsFiltered"`
	Data            []StockData `json:"data"`
}

// generateFetchStockURL --
func generateFetchStockURL(start int, length int) string {
	return fmt.Sprintf(
		`https://www.idx.co.id/umbraco/Surface/StockData/GetSecuritiesStock?start=%d&length=%d`,
		start,
		length,
	)
}

func (s *StockService) GetStockByCode(stockCode string) (*models.Stock, error) {
	stock := &models.Stock{}
	err := s.DB.Where("code = ?", stockCode).First(stock)
	if err != nil {
		return nil, err
	}
	return stock, nil
}

// FetchStocksFromDB -- Fetch Stocks Data from local DB
func (s *StockService) FetchStocksFromDB() ([]models.Stock, error) {
	stocks := []models.Stock{}
	err := s.DB.All(&stocks)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

// FetchStocks -- Fetch Stocks Data from IDX API
func (s *StockService) FetchStocks() ([]StockData, error) {
	start := 0
	length := 10
	URL := generateFetchStockURL(start, length)
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
		URL := generateFetchStockURL(start, length)
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

func (s *StockService) SaveStockToDB(code string, name string, listingDate string, shares int64, listingBoard nulls.String) error {
	stock := &models.Stock{
		Code:         code,
		Name:         name,
		ListingDate:  listingDate,
		Shares:       shares,
		ListingBoard: listingBoard,
	}
	err := s.DB.Save(stock)
	if err != nil {
		return err
	}
	return nil
}

// SaveStockDataToDB -- Receive Stock Data and save them all to database
func (s *StockService) SaveStockDataToDB(stocksData []StockData) error {
	for _, sd := range stocksData {
		stockModel := &models.Stock{}
		query := s.DB.Where("code = ?", sd.Code)
		err := query.First(stockModel)
		if err != nil {
			fmt.Println(err)
		}
		// If stock exists, update
		if stockModel.ID != uuid.Nil {
			stockModel.Shares = sd.Shares
			stockModel.ListingBoard = sd.ListingBoard
		} else { // Else Create
			stockModel = &models.Stock{
				Code:         sd.Code,
				Name:         sd.Name,
				ListingDate:  sd.ListingDate,
				Shares:       sd.Shares,
				ListingBoard: sd.ListingBoard,
			}
		}
		err = s.DB.Save(stockModel)
	}
	return nil
}
