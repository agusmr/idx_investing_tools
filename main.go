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
	"update-stock-list",
	"save-report",
	"calculate-ratios",
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

	// -- save-stock-data
	saveReportCommand := flag.NewFlagSet("save-report", flag.ExitOnError)
	// -- -- save-stock-data flags
	saveReportFilePathPtr := saveReportCommand.String(
		"file",
		"",
		"File path",
	)
	saveReportDirPathPtr := saveReportCommand.String(
		"dir",
		"",
		"Directory path",
	)

	// -- calculate-ratios
	calculateRatiosCommand := flag.NewFlagSet("calculate-ratios", flag.ExitOnError)
	// -- -- calculate-ratios flags
	calculateRatiosYearPtr := calculateRatiosCommand.Int(
		"year",
		0,
		"Which financial year to calculate the ratios. Currently available {2018|2019}",
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
		updateStockListCommand.Parse(os.Args) // Here args not used, but need to be passed in
	case "save-report":
		saveReportCommand.Parse(os.Args[2:])
	case "calculate-ratios":
		calculateRatiosCommand.Parse(os.Args[2:])
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

	// Handle save-report subCommand
	if saveReportCommand.Parsed() {
		if *saveReportFilePathPtr == "" && *saveReportDirPathPtr == "" {
			saveReportCommand.PrintDefaults()
			os.Exit(1)
		}
		if *saveReportFilePathPtr != "" && *saveReportDirPathPtr != "" {
			fmt.Println("Cannot use both -file and -dir at the same time")
			saveReportCommand.PrintDefaults()
			os.Exit(1)
		}
		if *saveReportFilePathPtr != "" {
			// Handle Save Stock Data File
			//fmt.Println(*saveReportFilePathPtr)
			filePath := *saveReportFilePathPtr
			err := handlers.SaveReportFile(filePath)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Save report finished")
			os.Exit(0)
		}
		if *saveReportDirPathPtr != "" {
			// Handle Save Stock Data directory
			//fmt.Println(*saveReportDirPathPtr)
			dirPath := *saveReportDirPathPtr
			fmt.Println(dirPath)
			err := handlers.SaveReportDir(dirPath)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Save report finished")
			os.Exit(0)
		}
	}

	// Handle calculate-ratios subcommand
	if calculateRatiosCommand.Parsed() {
		if *calculateRatiosYearPtr == 0 || *calculateRatiosYearPtr < 2018 {
			calculateRatiosCommand.PrintDefaults()
			os.Exit(1)
		}
		err := handlers.CalculateRatios(*calculateRatiosYearPtr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Calculate ratios finished. Ratios are saved to DB")
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
