package services

import (
	//"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ExcelReport --
type ExcelReport struct {
	File       *excelize.File
	Worksheets []string
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

func (exRep *ExcelReport) TotalAssets() float64 {
	sheetName := exRep.Worksheets[2] // The sheet name for statement of financial position
	totalAssetCell := exRep.File.SearchSheet(sheetName, "Total assets")[0]
	row := totalAssetCell[1:]
	currTotalAssetCell := "B" + row
	strVal := exRep.File.GetCellValue(sheetName, currTotalAssetCell)
	return excelFloatToFloat(strVal)
}

func (exRep *ExcelReport) TotalAssetsPrevious() float64 {
	sheetName := exRep.Worksheets[2]
	totalAssetCell := exRep.File.SearchSheet(sheetName, "Total assets")[0]
	row := totalAssetCell[1:]
	prevTotalAssetCell := "C" + row
	strVal := exRep.File.GetCellValue(sheetName, prevTotalAssetCell)
	return excelFloatToFloat(strVal)
}

func (exRep *ExcelReport) NetIncome() float64 {
	sheetName := exRep.Worksheets[3]
	totalProfitCell := exRep.File.SearchSheet(sheetName, "Total profit (loss)")[0]
	row := totalProfitCell[1:]
	netIncomeCell := "B" + row
	strVal := exRep.File.GetCellValue(sheetName, netIncomeCell)
	return excelFloatToFloat(strVal)
}

func (exRep *ExcelReport) ROA() float64 {
	averageAssets := (exRep.TotalAssets() + exRep.TotalAssetsPrevious()) / 2
	ROA := exRep.NetIncome() / averageAssets
	return ROA
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
	for index, _ := range worksheetMap {
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
