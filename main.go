package main

import (
	"flag"
	"fmt"
	"os"
	//"time"
	//"github.com/kevinjanada/idx_investing_tools/constants"
	//"github.com/kevinjanada/idx_investing_tools/services"
)

var commands = []string{
	"get-reports",
}

func main() {
	// SubCommands -------
	//  get-reports subcommand
	getReportsCommand := flag.NewFlagSet("get-reports", flag.ExitOnError)
	//  get-reports flags
	getReportsYearPtr := getReportsCommand.Int(
		"year",
		0,
		"Which financial year reports to download. Currently available {2017|2018|2019}",
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
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Handle get-reports subcommand
	if getReportsCommand.Parsed() {
		if *getReportsYearPtr == 0 {
			getReportsCommand.PrintDefaults()
			os.Exit(1)
		}
		fmt.Println("Download", *getReportsYearPtr)
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
