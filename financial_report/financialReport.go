package financial_report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// IdxJsonResponse -- Json response from IDX get financial report API
type IdxJSONResponse struct {
	Search      Search   `json:"Search"`
	ResultCount int      `json:"ResultCount"`
	Results     []Result `json:"Results"`
}

type Search struct {
	ReportType string `json:"ReportType"`
	KodeEmiten string `json:"KodeEmiten"`
	Year       string `json:"Year"`
	Periode    string `json:"Periode"`
	indexfrom  int    `json:"indexfrom"`
	pagesize   int    `json:"pagesize"`
}

type Result struct {
	KodeEmiten    string       `json:"KodeEmiten"`
	File_Modified string       `json:"File_Modified"`
	Report_Period string       `json:"Report_Period"`
	Report_Year   string       `json:"Report_Year"`
	NamaEmiten    string       `json:"NamaEmiten"`
	Attachments   []Attachment `json:"Attachments"`
}

type Attachment struct {
	Emiten_Code   string `json:"Emiten_Code"`
	File_ID       string `json:"File_ID"`
	File_Modified string `json:"File_Modified"`
	File_Name     string `json:"File_Name"`
	File_Path     string `json:"File_Path"`
	File_Size     int    `json:"File_Size"`
	File_Type     string `json:"File_Type"`
	Report_Period string `json:"Report_Period"`
	Report_Type   string `json:"Report_Type"`
	Report_Year   string `json:"Report_Year"`
	NamaEmiten    string `json:"NamaEmiten"`
}

func (self *IdxJSONResponse) Print() {
	res, err := json.MarshalIndent(self, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(res))
}

func (self *IdxJSONResponse) GetExcelReports() []Attachment {
	attachments := []Attachment{}
	for _, res := range self.Results {
		for _, att := range res.Attachments {
			if att.File_Type == ".xlsx" {
				attachments = append(attachments, att)
			}
		}
	}
	return attachments
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

func GetFinancialReports(year int, period int) *IdxJSONResponse {
	URL := GenerateURL(0, 1000, year, period)
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	jsonResponse := &IdxJSONResponse{}
	err = json.Unmarshal([]byte(body), jsonResponse)
	if err != nil {
		fmt.Println(err)
	}
	return jsonResponse
}
