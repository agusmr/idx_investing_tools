package main

import (
	"fmt"

	"github.com/kevinjanada/idx_investing_tools/services"
)

func main() {
	frService := services.FinancialReportService{}
	err := frService.FetchFinancialReports(2017, 3)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	err = frService.DownloadExcelReports()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
