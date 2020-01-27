package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/kevinjanada/idx_investing_tools/tools"
)

// FinancialReportService --
// 		Service for getting financial report data from IDX API which includes the links to download the financial reports
type FinancialReportService struct {
	*FinancialReportAPIResponse
}

// FetchFinancialReports --
// 		Fetch Financial Reports Data from IDX API.
//		Return FinancialReport struct for the selected year and period (trimester/triwulan)
func (frService *FinancialReportService) FetchFinancialReports(year int, period int) error {
	stockService, err := NewStockService("tools_development")
	if err != nil {
		return err
	}
	stocks, err := stockService.FetchStocksFromDB()
	if err != nil {
		return err
	}
	page := 1
	pageSize := 1

	aggregatedResponse := &FinancialReportAPIResponse{Year: strconv.Itoa(year), Period: fmt.Sprintf("trimester_%d", period)}
	for _, s := range stocks {
		URL := generateURL(page, pageSize, year, period, s.Code)
		fmt.Printf("Fetching data for %s %d trimester_%d\n", s.Code, year, period)
		resp, err := http.Get(URL)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		nextResponse := &FinancialReportAPIResponse{Year: strconv.Itoa(year), Period: fmt.Sprintf("trimester_%d", period)}
		err = tools.JSONToStruct(resp, nextResponse)
		if err != nil {
			return err
		}

		aggregatedResponse.Results = append(aggregatedResponse.Results, nextResponse.Results...)
	}
	frService.FinancialReportAPIResponse = aggregatedResponse
	return nil
}

// DownloadExcelReports --
// 		Download the excel reports for all stocks
func (frService *FinancialReportService) DownloadExcelReports() error {
	if frService.FinancialReportAPIResponse == nil {
		return errors.New("You need to FetchFinancialReports() first")
	}
	err := frService.FinancialReportAPIResponse.DownloadExcelReports()
	if err != nil {
		return err
	}
	return nil
}

func generateURL(page int, pageSize int, year int, period int, emitenCode string) string {
	return fmt.Sprintf(
		"https://www.idx.co.id/umbraco/Surface/ListedCompany/GetFinancialReport?indexFrom=%d&pageSize=%d&year=%d&reportType=rdf&periode=tw%d&kodeEmiten=%s",
		page,
		pageSize,
		year,
		period,
		emitenCode,
	)
}

func removeDuplicates(elements []Result) []Result {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []Result{}

	for _, el := range elements {
		fmt.Println(el.KodeEmiten)
		if encountered[el.KodeEmiten] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[el.KodeEmiten] = true
			// Append to result slice.
			result = append(result, el)
		}
	}
	// Return the new slice.
	return result
}

// FinancialReportAPIResponse -- Json response from IDX get financial report API
type FinancialReportAPIResponse struct {
	Year        string
	Period      string
	Search      Search   `json:"Search"`
	ResultCount int      `json:"ResultCount"`
	Results     []Result `json:"Results"`
}

// Search -- FRAPIResponseSearch field
type Search struct {
	ReportType string `json:"ReportType"`
	KodeEmiten string `json:"KodeEmiten"`
	Year       string `json:"Year"`
	Periode    string `json:"Periode"`
	Indexfrom  int    `json:"indexfrom"`
	Pagesize   int    `json:"pagesize"`
}

// Result -- FRAPIResponseSearch Results
type Result struct {
	KodeEmiten   string       `json:"KodeEmiten"`
	FileModified string       `json:"File_Modified"`
	ReportPeriod string       `json:"Report_Period"`
	ReportYear   string       `json:"Report_Year"`
	NamaEmiten   string       `json:"NamaEmiten"`
	Attachments  []Attachment `json:"Attachments"`
}

// Attachment -- FRAPIResponseAttachment object
type Attachment struct {
	EmitenCode   string `json:"Emiten_Code"`
	FileID       string `json:"File_ID"`
	FileModified string `json:"File_Modified"`
	FileName     string `json:"File_Name"`
	FilePath     string `json:"File_Path"`
	FileSize     int    `json:"File_Size"`
	FileType     string `json:"File_Type"`
	ReportPeriod string `json:"Report_Period"`
	ReportType   string `json:"Report_Type"`
	ReportYear   string `json:"Report_Year"`
	NamaEmiten   string `json:"NamaEmiten"`
}

// Print -- Pretty Print  struct
func (fr *FinancialReportAPIResponse) Print() {
	res, err := json.MarshalIndent(fr, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(res))
}

// GetExcelReportLinks -- Return attachments data of type excel
func (fr *FinancialReportAPIResponse) GetExcelAttachments() []Attachment {
	attachments := []Attachment{}
	for _, res := range fr.Results {
		for _, att := range res.Attachments {
			if att.FileType == ".xlsx" {
				attachments = append(attachments, att)
			}
		}
	}
	return attachments
}

// DownloadExcelReports -- Download all available excel reports
func (fr *FinancialReportAPIResponse) DownloadExcelReports() error {
	excelAttachments := fr.GetExcelAttachments()
	for _, att := range excelAttachments {
		directory := filepath.Join("files", "excel_reports", fr.Year, fr.Period)
		err := tools.Download(directory, att.FileName, att.FilePath)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
