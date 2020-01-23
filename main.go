package main

import (
	"fmt"
	"time"

	"github.com/kevinjanada/idx_investing_tools/services"
)

func main() {
	excelReports, err := services.OpenExcelFilesInDir("files/excel_reports/2019/trimester_3")
	if err != nil {
		fmt.Println(err)
	}

	statementService := services.NewStatementService()

	location, err := time.LoadLocation("Asia/Jakarta")
	for _, report := range excelReports {
		stockCode := report.EntityCode()
		stockTotalAssets := report.TotalAssets()
		stockNetIncome := report.NetIncome()
		date := time.Date(2019, 9, 30, 0, 0, 0, 0, location)

		// Insert Total Assets
		err = statementService.InsertUpdateStatementRow(
			stockCode,
			services.StatementFinancialPosition,
			services.RowTotalAssets,
			stockTotalAssets,
			date,
		)
		if err != nil {
			fmt.Println(err)
		}

		// Insert Net Income
		err = statementService.InsertUpdateStatementRow(
			stockCode,
			services.StatementProfitOrLoss,
			services.RowNetIncome,
			stockNetIncome,
			date,
		)
		if err != nil {
			fmt.Println(err)
		}
	}
}
