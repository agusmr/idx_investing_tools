package services

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gobuffalo/nulls"
)

type ExcelReportService struct {
	SheetNames map[string]int
}

func NewExcelReportService() *ExcelReportService {
	sheetNames := map[string]int{
		"general information": 1,
		"financial position":  2,
		"profit or loss":      3,
		"changes in equity":   4,
	}
	service := &ExcelReportService{SheetNames: sheetNames}
	return service
}

func (ex *ExcelReportService) LoadFile(filePath string) (*ExcelReport, error) {
	excelReport, err := NewExcelReport(filePath)
	if err != nil {
		return nil, err
	}
	return excelReport, nil
}

func (ex *ExcelReportService) LoadDir(dirPath string) ([]*ExcelReport, error) {
	excelReports, err := OpenExcelFilesInDir(dirPath)
	if err != nil {
		return nil, err
	}
	return excelReports, nil
}

// SaveReport
/**
  1. For now go to these sheets:
    - general information idx 1
    - financial position idx 2
    - profit or loss idx 3
  2. get date: cell B4
      Parse Date from 30 September 2017 to 30 Sep 17 00:00 GMT+7
      t, err := time.Parse(time.RFC822, "30 Sep 17 00:00 GMT+7")
    get stock_code: cell B7
  3. For each sheets:
    iterate through the rows (start from row 5 until row title is empty):
      get row title // Column D
      get the row amount // Current year amount is column B
      get row by title from db

      if row title !exists in db:
        insert title to db

      get row_fact from db
        where
          row_fact->date == date and
          row_fact->stock_code == stock_code
          row_fact->row_id == row
                              where row->row_title_id == title->id

      if row_fact exists:
        row_fact->amount = row_amount
        save(row_fact)
      else:
        row = new row {
          row_title_id = title->id
        }

        get the created row from db

        row_fact = new row_fact {
          row_id = row->id
          amount = amount
          date = date
        }
**/
func (ex *ExcelReportService) SaveReportToDB(exRep *ExcelReport) error {
	sheetNames := []string{
		"general information",
		"financial position",
		"profit or loss",
		"changes in equity",
	}
	statementService, err := NewStatementService("tools_development")
	if err != nil {
		return err
	}
	// Get Stock Code ----
	stockCode := exRep.EntityCode()
	// Get date -----
	date, err := exRep.Date()
	if err != nil {
		return err
	}
	// Check if Stock Code exists in DB ----
	stockService, err := NewStockService("tools_development")
	if err != nil {
		return err
	}
	_, err = stockService.GetStockByCode(stockCode)
	if err != nil {
		// if stock does not exist, add it to DB with limited information
		stockName := exRep.EntityName()
		err = stockService.SaveStockToDB(stockCode, stockName, "-", -1, nulls.NewString("-"))
		if err != nil {
			return err
		}
	}
	// -----
	fmt.Printf("Saving %s report to DB. Please wait..\n", stockCode)

	for _, sn := range sheetNames[1:3] {
		sheetIndex, ok := ex.SheetNames[sn]
		if !ok {
			return fmt.Errorf("sheet name: %s does not exist", sn)
		}
		// Iterate through sheet rows
		row := 5
		titleCol := "D"
		amountCol := "B"
		var title string
		for {
			// Get row title from excel report
			titleCell := fmt.Sprintf("%s%d", titleCol, row)
			title = exRep.GetContent(sheetIndex, titleCell)
			if title == "" {
				break
			}
			//if row title !exists in db:
			//insert title to db
			_ = statementService.InsertRowTitle(title)
			rowTitle, err := statementService.GetRowTitle(title)
			if err != nil {
				fmt.Println("is it here?")
				return err
			}

			// Get row amount from excel report
			amountCell := fmt.Sprintf("%s%d", amountCol, row)
			stringAmount := exRep.GetContent(sheetIndex, amountCell)
			floatAmount := excelFloatToFloat(stringAmount)

			err = statementService.InsertUpdateStatementRow(stockCode, sn, rowTitle.Title, floatAmount, date)
			if err != nil {
				return err
			}

			row++
		}
	}
	// Save Common Stocks end of period
	commonStocksTitle := "Common stocks - Equity position, end of the period"
	commonStocksValue := exRep.CommonStocks()
	_ = statementService.InsertRowTitle(commonStocksTitle)
	err = statementService.InsertUpdateStatementRow(
		stockCode,
		sheetNames[3],
		commonStocksTitle,
		commonStocksValue,
		date,
	)
	if err != nil {
		return err
	}

	// Insert EPS
	EPSTitle := "EPS"
	EPSValue := exRep.EPS()
	// Some reports do not have common stock amount
	// This results in -Infinity or Infinity EPS value
	// Workaround, set it to either large negative or positive amount
	if math.IsInf(EPSValue, -1) {
		EPSValue = -1000000000
	}
	if math.IsInf(EPSValue, 1) {
		EPSValue = 1000000000
	}
	statementService.InsertRowTitle(EPSTitle)
	err = statementService.InsertUpdateStatementRow(
		stockCode,
		"ratios",
		EPSTitle,
		EPSValue,
		date,
	)
	if err != nil {
		return err
	}

	return nil
}

// ExcelReport --
type ExcelReport struct {
	File       *excelize.File
	Worksheets []string
}

func (exRep *ExcelReport) GetContent(sheetIndex int, cell string) string {
	sheetName := exRep.Worksheets[sheetIndex]
	content := exRep.File.GetCellValue(sheetName, cell)
	return content
}

func (exRep *ExcelReport) EntityName() string {
	sheetName := exRep.Worksheets[1]
	entityName := exRep.File.GetCellValue(sheetName, "B5")
	return entityName
}

func (exRep *ExcelReport) EntityCode() string {
	sheetName := exRep.Worksheets[1]
	entityCode := exRep.File.GetCellValue(sheetName, "B7")
	return entityCode
}

func (exRep *ExcelReport) Date() (time.Time, error) {
	dateCell := "B4"
	dateString := exRep.GetContent(1, dateCell)
	if dateString == "" {
		return time.Time{}, fmt.Errorf("Date not found on report")
	}
	date, err := convertReportDateToTime(dateString)
	if err != nil || date.IsZero() {
		return time.Time{}, fmt.Errorf("Error in converting report date to time.Time")
	}
	return date, nil
}

func (exRep *ExcelReport) NetIncome() float64 {
	sheetName := exRep.Worksheets[3]
	totalProfitCell := exRep.File.SearchSheet(sheetName, "Total profit (loss)")[0]
	row := totalProfitCell[1:]
	netIncomeCell := "B" + row
	strVal := exRep.File.GetCellValue(sheetName, netIncomeCell)
	return excelFloatToFloat(strVal)
}

func (exRep *ExcelReport) CommonStocks() float64 {
	sheetName := exRep.Worksheets[4]

	eqPosEndOfPeriodCell := exRep.File.SearchSheet(sheetName, "Equity position, end of the period")[0]
	fmt.Println("Eq pos CELL: ", eqPosEndOfPeriodCell)
	// Find Digit
	re := regexp.MustCompile(`\d+`)
	eqPosEndOfPeriodRow := string(re.Find([]byte(eqPosEndOfPeriodCell)))
	fmt.Println("Eq Pos ROW: ", eqPosEndOfPeriodRow)

	commonStocksCell := exRep.File.SearchSheet(sheetName, "Common stocks")[0]
	fmt.Println("Common stocks CELL: ", commonStocksCell)
	// Find Alphabet
	re = regexp.MustCompile(`[A-Z]+`)
	commonStocksCol := string(re.Find([]byte(commonStocksCell)))
	fmt.Println("Common Stocks COL: ", commonStocksCol)

	endOfPeriodCommonStockcell := commonStocksCol + eqPosEndOfPeriodRow
	strVal := exRep.File.GetCellValue(sheetName, endOfPeriodCommonStockcell)
	fmt.Println("Common stocks strVal", strVal)
	return excelFloatToFloat(strVal)
}

func (exRep *ExcelReport) PreferredStock() float64 {
	sheetName := exRep.Worksheets[4]
	eqPosEndOfPeriodCell := exRep.File.SearchSheet(sheetName, "Equity position, end of the period")[0]
	eqPosEndOfPeriodRow := eqPosEndOfPeriodCell[1:]

	preferredStocksCell := exRep.File.SearchSheet(sheetName, "Preferred stocks")[0]
	preferredStocksCol := preferredStocksCell[:1]

	endOfPeriodPreferredStockValueCell := preferredStocksCol + eqPosEndOfPeriodRow
	strVal := exRep.File.GetCellValue(sheetName, endOfPeriodPreferredStockValueCell)
	return excelFloatToFloat(strVal)
}

func (exRep *ExcelReport) EPS() float64 {
	netIncome := exRep.NetIncome()
	preferredDividends := 0.0 // FIXME: Could not find this in reports yet
	commonStocks := exRep.CommonStocks()
	// net income * 1000 because the report values are in thousands
	EPS := ((netIncome * 1000) - preferredDividends) / commonStocks
	fmt.Println("Calculating EPS")
	fmt.Println("net income", netIncome)
	fmt.Println("common stocks", commonStocks)
	fmt.Println("net income * 1000 / common stocks = ", EPS)
	return EPS
}

// OpenExcelFilesInDir -- Open All the files in directory
func OpenExcelFilesInDir(directory string) ([]*ExcelReport, error) {
	var filepaths []string
	var files []*ExcelReport
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		filepaths = append(filepaths, path)
		//fmt.Printf("%+v\n", info)
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, path := range filepaths {
		fmt.Printf("Opening file %s\n", path)
		f, err := NewExcelReport(path)
		if err != nil {
			continue
		}
		files = append(files, f)
	}
	return files, nil
}

// NewExcelReport -- Open an excel file and instantiate ExcelReport struct
func NewExcelReport(filepath string) (*ExcelReport, error) {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, err
	}
	worksheets := ArrangeWorksheets(f.GetSheetMap())
	excelReport := &ExcelReport{File: f, Worksheets: worksheets}
	return excelReport, nil
}

// ArrangeWorksheets -- Arrange the worksheet names. sort in ascending order by index from excelize.File.GetSheetMap()
func ArrangeWorksheets(worksheetMap map[int]string) []string {
	var worksheets []string
	var indexes []int
	for index := range worksheetMap {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)
	for _, idx := range indexes {
		worksheets = append(worksheets, worksheetMap[idx])
	}
	return worksheets
}

func excelFloatToFloat(excelFloatValue string) float64 {
	split := strings.Split(excelFloatValue, "E")

	val, _ := strconv.ParseFloat(split[0], 64)
	if len(split) == 1 {
		return val
	}
	powerOf, _ := strconv.ParseFloat(split[1], 64)

	finalVal := val * math.Pow(10, powerOf)
	return finalVal
}

// prepend -- Helper to prepend to slice
func prepend(stringSlice []string, elem string) []string {
	stringSlice = append(stringSlice, "")
	copy(stringSlice[1:], stringSlice)
	stringSlice[0] = elem
	return stringSlice
}

// insert -- Helper to insert to slice by index
func insert(strSlice []string, elem string, i int) []string {
	strSlice = append(strSlice[:i], append([]string{elem}, strSlice[i:]...)...)
	return strSlice
}

// convertReportDateToTime -- Helper to convert date in excel report to time.Time
func convertReportDateToTime(dateString string) (time.Time, error) {
	if dateString[0] == '\'' {
		dateString = dateString[1:]
	}
	dateSplit := strings.Split(dateString, " ")
	day := dateSplit[0]
	month := dateSplit[1]
	year := dateSplit[2]

	month = month[:3]
	year = year[2:]

	finalDateString := day + " " + month + " " + year + " 00:00 GMT+7"
	t, err := time.Parse(time.RFC822, finalDateString)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
