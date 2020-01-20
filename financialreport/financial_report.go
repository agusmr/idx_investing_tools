package financialreport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/kevinjanada/idx_investing_tools/tools"
)

// FinancialReports -- Json response from IDX get financial report API
type FinancialReports struct {
	Year        string
	Period      string
	Search      Search   `json:"Search"`
	ResultCount int      `json:"ResultCount"`
	Results     []Result `json:"Results"`
}

// Search -- FinancialReports Search field
type Search struct {
	ReportType string `json:"ReportType"`
	KodeEmiten string `json:"KodeEmiten"`
	Year       string `json:"Year"`
	Periode    string `json:"Periode"`
	Indexfrom  int    `json:"indexfrom"`
	Pagesize   int    `json:"pagesize"`
}

// Result -- FinancialReports Search Results
type Result struct {
	KodeEmiten   string       `json:"KodeEmiten"`
	FileModified string       `json:"File_Modified"`
	ReportPeriod string       `json:"Report_Period"`
	ReportYear   string       `json:"Report_Year"`
	NamaEmiten   string       `json:"NamaEmiten"`
	Attachments  []Attachment `json:"Attachments"`
}

// Attachment -- FinancialReports Attachment object
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
func (fr *FinancialReports) Print() {
	res, err := json.MarshalIndent(fr, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(res))
}

// GetExcelReportLinks -- Return download links of all excel reports
func (fr *FinancialReports) GetExcelReportLinks() []Attachment {
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
func (fr *FinancialReports) DownloadExcelReports() error {
	excelReportLinks := fr.GetExcelReportLinks()
	for _, report := range excelReportLinks {
		directory := filepath.Join("files", "excel_reports", fr.Year, fr.Period)
		err := tools.Download(directory, report.FileName, report.FilePath)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

// GenerateURL --
func GenerateURL(page int, pageSize int, year int, period int) string {
	return fmt.Sprintf(
		"https://www.idx.co.id/umbraco/Surface/ListedCompany/GetFinancialReport?indexFrom=%d&pageSize=%d&year=%d&reportType=rdf&periode=tw%d&kodeEmiten=",
		page,
		pageSize,
		year,
		period,
	)
}

// GetFinancialReports -- Return FinancialReport struct for the selected year and period (trimester/triwulan)
func GetFinancialReports(year int, period int) *FinancialReports {
	URL := GenerateURL(0, 1000, year, period)
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	financialReports := &FinancialReports{Year: strconv.Itoa(year), Period: fmt.Sprintf("trimester_%d", period)}
	err = json.Unmarshal([]byte(body), financialReports)
	if err != nil {
		fmt.Println(err)
	}
	return financialReports
}
