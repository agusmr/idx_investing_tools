package main

import ( // "time"
	"fmt"
	"time"

	"github.com/kevinjanada/idx_investing_tools/constants"
	"github.com/kevinjanada/idx_investing_tools/services"
)

func main() {
	// statementService := services.NewStatementService("api_development")
	// err := statementService.NewStatement(constants.StatementChangesInEquity)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// testFile, err := services.NewExcelReport("files/excel_reports/2017/trimester_3/FinancialStatement-2017-III-AALI.xlsx")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// preferredStockValue := testFile.PreferredStock()
	// fmt.Println(preferredStockValue)

	excelReports, err := services.OpenExcelFilesInDir("files/excel_reports/2017/trimester_3")
	if err != nil {
		fmt.Println(err)
	}

	statementService := services.NewStatementService("api_development")
	// rowTitleID, err := uuid.FromBytes([]byte("18374f1e-d432-4058-8597-a8ef471eecdf"))
	// stockID, err := uuid.FromBytes([]byte("b227b28b-3139-458a-b38d-11c64ec1f9f1"))
	// date, err := time.Parse(time.RFC3339, "2017-09-30T00:00:00Z")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// statementRow := statementService.StatementRowExists(rowTitleID, stockID, date)
	// fmt.Printf("%+v", statementRow)
	// err = statementService.NewStatementRowTitle(constants.RowPreferredStocks)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	location, err := time.LoadLocation("Asia/Jakarta")
	for _, report := range excelReports {
		stockCode := report.EntityCode()
		stockPreferredStock := report.PreferredStock()
		fmt.Printf("%s preferred-stock: %f\n", stockCode, stockPreferredStock)
		date := time.Date(2017, 9, 30, 0, 0, 0, 0, location)

		// Insert Preferred Stock
		err = statementService.InsertUpdateStatementRow(
			stockCode,
			constants.StatementChangesInEquity,
			constants.RowPreferredStocks,
			stockPreferredStock,
			date,
		)
		if err != nil {
			fmt.Println(err)
		}

		// stockTotalAssets := report.TotalAssets()
		// // Insert Total Assets
		// err = statementService.InsertUpdateStatementRow(
		// 	stockCode,
		// 	constants.StatementFinancialPosition,
		// 	constants.RowTotalAssets,
		// 	stockTotalAssets,
		// 	date,
		// )
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// stockNetIncome := report.NetIncome()
		// // Insert Net Income
		// err = statementService.InsertUpdateStatementRow(
		// 	stockCode,
		// 	constants.StatementProfitOrLoss,
		// 	constants.RowNetIncome,
		// 	stockNetIncome,
		// 	date,
		// )
		// if err != nil {
		// 	fmt.Println(err)
		// }
	}
}
