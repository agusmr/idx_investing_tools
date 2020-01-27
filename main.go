package main

import (
	"flag"
	"fmt"
	"os"
	//"time"
	//"github.com/kevinjanada/idx_investing_tools/constants"
	//"github.com/kevinjanada/idx_investing_tools/services"
	"github.com/kevinjanada/idx_investing_tools/handlers"
)

var commands = []string{
	"get-reports",
}

func main() {
	// SubCommands -------
	//  -- update-stock-list
	updateStockListCommand := flag.NewFlagSet("update-stock-list", flag.ExitOnError)

	//  -- get-reports
	getReportsCommand := flag.NewFlagSet("get-reports", flag.ExitOnError)
	//  -- -- get-reports flags
	getReportsYearPtr := getReportsCommand.Int(
		"year",
		0,
		"Which financial year reports to download. Currently available {2017|2018|2019}",
	)
	getReportsPeriodPtr := getReportsCommand.Int(
		"period",
		0,
		"Period of the year in trimester {1|2|3}",
	)

	// No Subcommands error handler
	if len(os.Args) < 2 {
		//fmt.Println("")
		fmt.Println("Need to input commands")
		fmt.Println("   Available commands:")
		for _, c := range commands {
			fmt.Printf("          %s\n", c)
		}
		os.Exit(1)
	}

	// Parse Subcommands
	switch os.Args[1] {
	case "get-reports":
		getReportsCommand.Parse(os.Args[2:])
	case "update-stock-list":
		updateStockListCommand.Parse(os.Args)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Handle update-stock-list subcommand
	if updateStockListCommand.Parsed() {
		err := handlers.UpdateStockList()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Stocks updated")
		os.Exit(0)
	}

	// Handle get-reports subcommand
	if getReportsCommand.Parsed() {
		if *getReportsYearPtr == 0 {
			getReportsCommand.PrintDefaults()
			os.Exit(1)
		}
		if *getReportsPeriodPtr == 0 {
			getReportsCommand.PrintDefaults()
			os.Exit(1)
		}
		year := *getReportsYearPtr
		period := *getReportsPeriodPtr
		err := handlers.GetReports(year, period)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Download reports finished")
		os.Exit(0)
	}

	//excelReports, err := services.OpenExcelFilesInDir("files/excel_reports/2017/trimester_3")
	//if err != nil {
	//fmt.Println(err)
	//}

	//statementService := services.NewStatementService("api_development")

	//location, err := time.LoadLocation("Asia/Jakarta")
	//for _, report := range excelReports {
	//stockCode := report.EntityCode()
	//stockPreferredStock := report.PreferredStock()
	//fmt.Printf("%s preferred-stock: %f\n", stockCode, stockPreferredStock)
	//date := time.Date(2017, 9, 30, 0, 0, 0, 0, location)

	//// Insert Preferred Stock
	//err = statementService.InsertUpdateStatementRow(
	//stockCode,
	//constants.StatementChangesInEquity,
	//constants.RowPreferredStocks,
	//stockPreferredStock,
	//date,
	//)
	//if err != nil {
	//fmt.Println(err)
	//}
	//}
}
