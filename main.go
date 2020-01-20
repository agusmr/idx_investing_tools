package main

import (
	"fmt"

	fr "github.com/kevinjanada/idx_investing_tools/financialreport"
)

func main() {
	financialReports := fr.GetFinancialReports(2019, 3)
	err := financialReports.DownloadExcelReports()
	if err != nil {
		fmt.Println(err)
	}
}
