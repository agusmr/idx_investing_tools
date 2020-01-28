package handlers

import (
	"fmt"
	"github.com/kevinjanada/idx_investing_tools/services"
)

func SaveReportFile(filePath string) error {
	fmt.Println(filePath)
	exRepService := services.NewExcelReportService()
	exFile, err := exRepService.LoadFile(filePath)
	if err != nil {
		return err
	}
	err = exRepService.SaveReportToDB(exFile)
	return nil
}

func SaveReportDir(dir string) error {
	// TODO:
	return nil
}
