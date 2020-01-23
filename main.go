package main

import (
	"fmt"

	fr "github.com/kevinjanada/idx_investing_tools/financialreport"
)

func main() {
	stocks, _ := fr.FetchStocksFromDB()
	fmt.Printf("%+v", stocks)
}
