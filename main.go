package main

import (
	"fmt"
	fr "github.com/kevinjanada/idx_investing_tools/financialreport"
	//"math"
	//"strconv"
	//"strings"
)

func main() {
	//financialReports := fr.GetFinancialReports(2019, 3)
	//err := financialReports.DownloadExcelReports()
	//if err != nil {
	//fmt.Println(err)
	//}

	//filepath := "files/excel_reports/2019/trimester_3/FinancialStatement-2019-III-AALI.xlsx"
	//testFile, err := fr.OpenExcelFile(filepath)
	//if err != nil {
	//fmt.Println(err)
	//}

	dir := "files/excel_reports/2019/trimester_3"

	files, err := fr.OpenExcelFilesInDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		//fmt.Println(file.Worksheets)

		entityName := file.EntityName()
		entityCode := file.EntityCode()
		totalAssetsCurrent := file.TotalAssets()
		totalAssetsPrevious := file.TotalAssetsPrevious()
		netIncome := file.NetIncome()
		ROA := file.ROA()

		if ROA > 0.24 {
			fmt.Printf("Entity Name: %s \n", entityName)
			fmt.Printf("Entity Code: %s \n", entityCode)
			fmt.Printf("Current total assets %f \n", totalAssetsCurrent)
			fmt.Printf("Previous total assets %f \n", totalAssetsPrevious)
			fmt.Printf("Net Income %f \n", netIncome)
			fmt.Printf("ROA %f \n", ROA)

			fmt.Println("==================================")
			fmt.Println("----------------------------------")
			fmt.Println()
		}

	}

}
