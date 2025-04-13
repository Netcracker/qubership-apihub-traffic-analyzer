// Copyright 2024-2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
)

var (
	defOptions = excelize.Options{
		ShortDatePattern: "dd-mmm",
		LongDatePattern:  "yyyy-mm-dd;@",
		LongTimePattern:  "hh:mm:ss",
	}
	NoExcelFileToWrite = errors.New("no excel file to write")
)

type ExcelService interface {
	CreateFile(filename string, bSave bool) error
	SetCellsValues(sheetName string, columnsValue map[string]interface{}) error
	SetCellValue(sheetName string, cellAddr string, value interface{}) error
	MakeNewSheet(sheetName string) error
	RemoveSheet(sheetName string) error
	CloseFile() error
	SetFilter(sheetName, cellsRange string) error
	WriteTo(w io.Writer) (int64, error)
	SetColumnWidth(sheetName string, begin, end string, width float64) error
}

type excelService struct {
	report *excelize.File
}

func newXlFile(fileName string, bSave bool) (*excelize.File, error) {
	excelReportFile := excelize.NewFile(defOptions)
	if excelReportFile == nil {
		return excelReportFile, fmt.Errorf("unable to create new Excel file")
	}
	excelReportFile.Path = fileName
	if bSave {
		err := excelReportFile.Save(defOptions)
		if err != nil {
			return excelReportFile, err
		}
	}
	return excelReportFile, nil
}

func NewExcelService(fileName string) (ExcelService, error) {
	excelReportFile, err := newXlFile(fileName, false)
	if err != nil {
		return nil, err
	}
	return &excelService{report: excelReportFile}, nil
}

func (xl *excelService) SetCellsValues(sheetName string, columnsValue map[string]interface{}) error {
	if xl.report == nil {
		return NoExcelFileToWrite
	}
	for key, value := range columnsValue {
		err := xl.report.SetCellValue(sheetName, key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (xl *excelService) SetCellValue(sheetName string, cellAddr string, value interface{}) error {
	if xl.report == nil {
		return NoExcelFileToWrite
	}
	err := xl.report.SetCellValue(sheetName, cellAddr, value)
	if err != nil {
		return err
	}
	return nil
}

func (xl *excelService) CloseFile() error {
	if xl.report == nil {
		return nil
	}
	err := xl.report.Save(defOptions)
	if err != nil {
		return err
	}
	err = xl.report.Close()
	if err == nil {
		xl.report = nil
	}
	return err
}

func (xl *excelService) CreateFile(fileName string, bSave bool) error {
	if xl.report != nil {
		return fmt.Errorf("unable to create %s : Excel file is already created", fileName)
	}
	excelReportFile, err := newXlFile(fileName, bSave)
	if err != nil {
		return err
	}
	xl.report = excelReportFile
	return nil
}

func (xl *excelService) MakeNewSheet(sheetName string) error {
	if xl.report == nil {
		return NoExcelFileToWrite
	}
	_, err := xl.report.NewSheet(sheetName)
	return err
}

func (xl *excelService) RemoveSheet(sheetName string) error {
	if xl.report == nil {
		return NoExcelFileToWrite
	}

	return xl.report.DeleteSheet(sheetName)
}

func (xl *excelService) SetFilter(sheetName, cellsRange string) error {
	return xl.report.AddTable(sheetName, &excelize.Table{Range: cellsRange})
}

func (xl *excelService) WriteTo(w io.Writer) (int64, error) {
	err := xl.report.Save(defOptions)
	if err == nil {
		err = xl.report.Write(w)
	}
	return 0, err
}

func (xl *excelService) SetColumnWidth(sheetName string, begin, end string, width float64) error {
	return xl.report.SetColWidth(sheetName, begin, end, width)
}
