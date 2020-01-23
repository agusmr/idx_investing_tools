package main

import (
	"fmt"

	fr "github.com/kevinjanada/idx_investing_tools/financialreport"
)

func main() {
	frReportResponses, err := fr.GetFinancialReports(2018, 3)
	if err != nil {
		fmt.Println(err)
	}

	frReportResponses.DownloadExcelReports()
}
