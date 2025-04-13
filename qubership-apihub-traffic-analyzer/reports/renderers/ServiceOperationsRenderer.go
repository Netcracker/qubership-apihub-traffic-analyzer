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

package renderers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/reports/generators"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/repository"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/service"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/utils"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
)

type ServiceOperationsRenderer struct {
	db             db.ConnectionProvider
	req            view.ReportDataRequest
	report         entities.ReportEntity
	reportType     entities.ReportTypeEntity
	xl             service.ExcelService
	currentDataRow int
	maxColIdx      int
	fileName       string
	reportFile     *os.File
	serviceName    string
}

const (
	paramsSheetIndex        = 0
	dataSheetIndex          = 1
	letterA                 = 65
	unsupportedRenderFormat = "report format %s is not currently supported"
	nilFile                 = "%s file %s not initialized"
	jsonObjectBegin         = "{"
	jsonObjectEnd           = "}"
	jsonArrayBegin          = "["
	jsonArrayEnd            = "]"
	jsonArraySep            = ","
	defaultSheetName        = "Sheet1"
)

var (
	sheets        = []string{"Parameters", "Data"}
	columnHeaders = []string{"Sender", "Receiver", "Method", "Path", "Operation-id", "Count", "Comment"}
	colWidths     = []float64{20., 20., 10., 75., 70., 10., 15.}

	//byteArrayBegin  = []byte(jsonArrayBegin)
	//byteArrayEnd    = []byte(jsonArrayEnd)
	byteArraySep = []byte(jsonArraySep)
	//byteObjectBegin = []byte(jsonObjectBegin)
	//byteObjectEnd   = []byte(jsonObjectEnd)
)

// NewServiceOperationsRenderer
// created a renderer for service operation report
func NewServiceOperationsRenderer(db db.ConnectionProvider,
	rqi interface{},
	workDir string) (ReportRenderer, error) {
	var (
		xl         service.ExcelService = nil
		err        error                = nil
		fileName   string
		fh         *os.File = nil
		report     *entities.ReportEntity
		reportType *entities.ReportTypeEntity
	)
	req := rqi.(view.ReportDataRequest)
	// check report existence and type
	report, reportType, err = repository.GetReport(db, req.Id, string(generators.ServiceOperationReport))
	if err != nil {
		return nil, err
	}
	switch req.Format {
	case view.ReportFormatExcel:
		{
			fileName = path.Join(workDir, utils.MakeUniqueId()+reportType.Name+view.ReportFileExtExcel)
			xl, err = service.NewExcelService(fileName)
			if err != nil {
				return nil, fmt.Errorf("unable to create excel file %s: %v", fileName, err)
			}
		}
	case view.ReportFormatJson:
		{
			fileName = reportType.Name + view.ReportFileExtJson
			fh, err = os.CreateTemp(workDir, fileName)
			if err != nil {
				return nil, fmt.Errorf("unable to create intermediate file %s: %v", fileName, err)
			}
		}
	default:
		return nil, fmt.Errorf(unsupportedRenderFormat, req.Format)
	}
	return &ServiceOperationsRenderer{
		db:             db,
		req:            req,
		report:         *report,
		reportType:     *reportType,
		xl:             xl,
		currentDataRow: 2,
		maxColIdx:      0,
		fileName:       fileName,
		reportFile:     fh,
	}, nil
}

func (srr *ServiceOperationsRenderer) MakeReportHeader() error {
	switch srr.req.Format {
	case view.ReportFormatExcel:
		{
			if srr.xl == nil {
				return service.NoExcelFileToWrite
			}
			for _, sheet := range sheets {
				err := srr.xl.MakeNewSheet(sheet)
				if err != nil {
					log.Debugf("unable to create sheet %s: %v", sheet, err)
				}
			}
			err := srr.xl.RemoveSheet(defaultSheetName)
			if err != nil {
				log.Debugf("unable to delete default sheet %v", err)
			}
			colValues := make(map[string]interface{})
			rowNum := 1
			colValues[fmt.Sprintf("A%d", rowNum)] = "Report type:"
			colValues[fmt.Sprintf("B%d", rowNum)] = srr.reportType.Name
			rowNum++
			params, err := view.UnmarshalServiceReportRequest([]byte(srr.report.ReportParameters))
			if err != nil {
				return fmt.Errorf("unable to unmarshall report parameters: %v", err)
			}
			colValues[fmt.Sprintf("A%d", rowNum)] = "Capture Id:"
			colValues[fmt.Sprintf("B%d", rowNum)] = params.CaptureId
			rowNum++
			colValues[fmt.Sprintf("A%d", rowNum)] = "Service name:"
			colValues[fmt.Sprintf("B%d", rowNum)] = params.ServiceName
			srr.serviceName = params.ServiceName
			rowNum++
			colValues[fmt.Sprintf("A%d", rowNum)] = "Service version:"
			colValues[fmt.Sprintf("B%d", rowNum)] = params.ServiceVersion
			rowNum++
			colValues[fmt.Sprintf("A%d", rowNum)] = "Version status:"
			colValues[fmt.Sprintf("B%d", rowNum)] = params.VersionStatus
			rowNum++
			colValues[fmt.Sprintf("A%d", rowNum)] = "Requested at:"
			colValues[fmt.Sprintf("B%d", rowNum)] = srr.report.CreatedAt
			rowNum++
			colValues[fmt.Sprintf("A%d", rowNum)] = "Completed at:"
			colValues[fmt.Sprintf("B%d", rowNum)] = srr.report.CompletedAt
			err = srr.xl.SetCellsValues(sheets[paramsSheetIndex], colValues)
			if err != nil {
				return fmt.Errorf("unable to fill %s sheet : %v", sheets[paramsSheetIndex], err)
			}
			err = srr.xl.SetColumnWidth(sheets[paramsSheetIndex], "A", "A", 20.)
			err = srr.xl.SetColumnWidth(sheets[paramsSheetIndex], "B", "B", 40.)

			colHeadValues := make(map[string]interface{})
			for i, header := range columnHeaders {
				colLetter := string(byte(i + letterA))
				srr.maxColIdx = i
				colHeadValues[fmt.Sprintf("%s1", colLetter)] = header
				err = srr.xl.SetColumnWidth(sheets[dataSheetIndex], colLetter, colLetter, colWidths[i])
				if err != nil {
					log.Debugf("unable to set column %s:%s width for header %s: %v", sheets[dataSheetIndex], colLetter, header, err)
				}
			}
			return srr.xl.SetCellsValues(sheets[dataSheetIndex], colHeadValues)
		}
	case view.ReportFormatJson:
		{
			_, err := srr.reportFile.WriteString(fmt.Sprintf("%s\n\"parameters\":", jsonObjectBegin))
			if err != nil {
				return fmt.Errorf("unable to write object begin: %v", err)
			}
			_, err = srr.reportFile.Write([]byte(srr.report.ReportParameters))
			if err != nil {
				return fmt.Errorf("unable to write report parameters: %v", err)
			}
			_, err = srr.reportFile.WriteString(fmt.Sprintf(", \"data\":%s\n", jsonArrayBegin))
			if err != nil {
				return fmt.Errorf("unable to write array begin: %v", err)
			}
		}
	}
	return nil
}

// ProcessRows
// iterates all the data rows and render e
func (srr *ServiceOperationsRenderer) ProcessRows() error {
	reportData := make([]entities.ReportDataRow, 0)
	reportQuery := srr.db.GetConnection().Model(&reportData).Where("report_id=?", srr.report.ReportId)
	if reportQuery == nil {
		return fmt.Errorf("no report data found for report id %d", srr.report.ReportId)
	}
	err := reportQuery.Select()
	if err != nil {
		if !errors.Is(err, pg.ErrNoRows) {
			return err
		}
	}
	for _, reportDataRow := range reportData {
		err = srr.RenderRow(&reportDataRow)
		if err != nil {
			break
		}
	}
	return err
}

// RenderRow
// render data row into report
func (srr *ServiceOperationsRenderer) RenderRow(dataRow *entities.ReportDataRow) error {
	data, err := view.DecodeOperationStatusWithPeers([]byte(dataRow.ReportRow))
	if srr.req.Format == view.ReportFormatExcel {
		if err != nil {
			return err
		}
		colValues := make(map[string]interface{})
		if len(data.Source) > 2 {
			colValues[fmt.Sprintf("A%d", srr.currentDataRow)] = data.Source
		}
		if len(data.Destination) > 2 {
			colValues[fmt.Sprintf("B%d", srr.currentDataRow)] = data.Destination
		} else {
			colValues[fmt.Sprintf("B%d", srr.currentDataRow)] = srr.serviceName
		}
		colValues[fmt.Sprintf("C%d", srr.currentDataRow)] = data.Method
		colValues[fmt.Sprintf("D%d", srr.currentDataRow)] = data.Path
		colValues[fmt.Sprintf("E%d", srr.currentDataRow)] = data.Id
		colValues[fmt.Sprintf("F%d", srr.currentDataRow)] = data.HitCount
		colValues[fmt.Sprintf("G%d", srr.currentDataRow)] = data.Status
		err = srr.xl.SetCellsValues(sheets[dataSheetIndex], colValues)
		if err == nil {
			srr.currentDataRow++
		}
		return err
	} else {
		if srr.currentDataRow > 2 {
			_, err = srr.reportFile.Write(byteArraySep)
			if err != nil {
				log.Debugf("unable to write JSON array separator: %v", err)
			}
		}
		_, err = srr.reportFile.Write([]byte(dataRow.ReportRow))
		if err != nil {
			log.Debugf("unable to write JSON array element: %v", err)
		}
		srr.currentDataRow++
	}
	return nil
}

// MakeReportFooter
// write format dependent footer
func (srr *ServiceOperationsRenderer) MakeReportFooter() error {
	if srr.req.Format == view.ReportFormatExcel {
		return srr.xl.SetFilter(sheets[dataSheetIndex],
			fmt.Sprintf("A1:%s%d", string(byte(srr.maxColIdx+letterA)), srr.currentDataRow))
	} else {
		_, err := srr.reportFile.WriteString(fmt.Sprintf("%s%s", jsonArrayEnd, jsonObjectEnd))
		if err != nil {
			log.Debugf("unable to write JSON end: %v", err)
		}
	}
	return nil
}

// FlushData
// flush report bytes into
func (srr *ServiceOperationsRenderer) FlushData(w io.Writer) error {
	switch srr.req.Format {
	case view.ReportFormatExcel:
		{
			// copy binary stream to channel
			_, err := srr.xl.WriteTo(w)
			return err
		}
	case view.ReportFormatJson:
		{
			_, err := srr.reportFile.Seek(0, io.SeekStart)
			if err == nil {
				_, err = srr.reportFile.WriteTo(w)
			}
			return err
		}
	}
	return fmt.Errorf(unsupportedRenderFormat, srr.req.Format)
}

// GetFileName
// returns file name without path for HTTP header
func (srr *ServiceOperationsRenderer) GetFileName() string {
	return path.Base(srr.fileName)
}

// Dispose
// try to dispose internally allocated resources
func (srr *ServiceOperationsRenderer) Dispose() {
	var err error = nil
	if srr.fileName == view.EmptyString {
		log.Debugf("nothing to perform for empty file name")
		return
	}
	switch srr.req.Format {
	case view.ReportFormatExcel:
		if srr.xl != nil {
			err = srr.xl.CloseFile()
		} else {
			err = fmt.Errorf(nilFile, srr.req.Format, srr.fileName)
		}
	case view.ReportFormatJson:
		if srr.reportFile != nil {
			err = srr.reportFile.Close()
		} else {
			err = fmt.Errorf(nilFile, srr.req.Format, srr.fileName)
		}
	default:
		err = fmt.Errorf(unsupportedRenderFormat, srr.req.Format)
	}
	if err != nil {
		log.Debugf("unable to close file %s: %v", srr.fileName, err)
	}
	err = os.Remove(srr.fileName)
	if err != nil {
		log.Debugf("unable to delete file %s: %v", srr.fileName, err)
	}
}
