package main

import (
	"fmt"
	"time"

	"github.com/kevinjanada/idx_investing_tools/constants"
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
			constants.StatementFinancialPosition,
			constants.RowTotalAssets,
			stockTotalAssets,
			date,
		)
		if err != nil {
			fmt.Println(err)
		}

		// Insert Net Income
		err = statementService.InsertUpdateStatementRow(
			stockCode,
			constants.StatementProfitOrLoss,
			constants.RowNetIncome,
			stockNetIncome,
			date,
		)
		if err != nil {
			fmt.Println(err)
		}
	}
}
