package main

import (
	"fmt"
	fr "github.com/kevinjanada/idx_investing_tools/financial_report"
	"github.com/kevinjanada/idx_investing_tools/tools"
)

func main() {
	result := fr.GetFinancialReports(2019, 3)
	excelReports := result.GetExcelReports()

	for _, report := range excelReports {
		filepath := fmt.Sprintf("files/excel_reports/%s", report.File_Name)

		err := tools.Download(filepath, report.File_Path)
		if err != nil {
			fmt.Println(err)
		}
	}
}
