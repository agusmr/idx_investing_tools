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
	if err != nil {
		return err
	}
	return nil
}

func SaveReportDir(dir string) error {
	exRepService := services.NewExcelReportService()
	exFiles, err := exRepService.LoadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range exFiles {
		err = exRepService.SaveReportToDB(f)
		if err != nil {
			return err
		}
	}
	return nil
}
