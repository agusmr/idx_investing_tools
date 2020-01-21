package financialreport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/kevinjanada/idx_investing_tools/tools"
)

// FRAPIResponse -- Json response from IDX get financial report API
type FRAPIResponse struct {
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
func (fr *FRAPIResponse) Print() {
	res, err := json.MarshalIndent(fr, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(res))
}

// GetExcelReportLinks -- Return download links of all excel reports
func (fr *FRAPIResponse) GetExcelReportLinks() []Attachment {
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
func (fr *FRAPIResponse) DownloadExcelReports() error {
	excelReportLinks := fr.GetExcelReportLinks()
	for _, report := range excelReportLinks {
		directory := filepath.Join("files", "excel_reports", fr.Year, fr.Period)
		fmt.Printf("Downloading %s", report.FileName)
		err := tools.Download(directory, report.FileName, report.FilePath)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

// GenerateURL --
func GenerateURL(page int, pageSize int, year int, period int, emitenCode string) string {
	return fmt.Sprintf(
		"https://www.idx.co.id/umbraco/Surface/ListedCompany/GetFinancialReport?indexFrom=%d&pageSize=%d&year=%d&reportType=rdf&periode=tw%d&kodeEmiten=%s",
		page,
		pageSize,
		year,
		period,
		emitenCode,
	)
}

// GetFinancialReports -- Return FinancialReport struct for the selected year and period (trimester/triwulan)
func GetFinancialReports(year int, period int) (*FRAPIResponse, error) {
	page := 1
	pageSize := 50
	URL := GenerateURL(page, pageSize, year, period, "")

	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	aggregatedResponse := &FRAPIResponse{Year: strconv.Itoa(year), Period: fmt.Sprintf("trimester_%d", period)}
	tools.JSONToStruct(resp, aggregatedResponse)
	if err != nil {
		return nil, err
	}

	resultCount := aggregatedResponse.ResultCount
	resultLeft := resultCount
	// While there are still result left
	for resultLeft > 0 {
		URL = GenerateURL(page, pageSize, year, period, "")
		fmt.Println(URL)
		resp, err := http.Get(URL)
		if err != nil {
			return nil, err
		}
		// While response is not 200
		for resp.StatusCode != 200 {
			// Try less pageSize
			pageSize--
			URL = GenerateURL(page, pageSize, year, period, "")
			fmt.Println(URL)
			resp, err = http.Get(URL)
			if err != nil {
				return nil, err
			}
		}
		defer resp.Body.Close()
		nextResponse := &FRAPIResponse{Year: strconv.Itoa(year), Period: fmt.Sprintf("trimester_%d", period)}
		err = tools.JSONToStruct(resp, nextResponse)
		if err != nil {
			return nil, err
		}

		aggregatedResponse.Results = append(aggregatedResponse.Results, nextResponse.Results...)
		page++
		resultLeft -= pageSize
		// Return pageSize to 50
		pageSize = 50
	}

	// Filter aggregatedResponse.Results
	aggregatedResponse.Results = removeDuplicates(aggregatedResponse.Results)

	return aggregatedResponse, nil
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
