package handlers

import (
	//"fmt"
	"github.com/kevinjanada/idx_investing_tools/services"
)

// GetReports -- params: year int, period int
func GetReports(year int, period int) error {
	frService := &services.FinancialReportService{}
	err := frService.FetchFinancialReports(year, period)
	if err != nil {
		return err
	}
	err = frService.DownloadExcelReports()
	if err != nil {
		return err
	}
	return nil
}
