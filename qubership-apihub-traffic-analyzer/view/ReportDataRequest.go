package view

import (
	"errors"
	"fmt"
)

const (
	ReportFormatJson  = "json"
	ReportFormatHtml  = "html"
	ReportFormatXml   = "xml"
	ReportFormatExcel = "excel"
	// ReportFileExtDot any file extension begins with
	ReportFileExtDot = "."
	// report file extensions

	// ReportFileExtJson JSON file format
	ReportFileExtJson = ReportFileExtDot + ReportFormatJson

	//ReportFileExtHtml = ReportFileExtDot + ReportFormatHtml
	//ReportFileExtXml  = ReportFileExtDot + ReportFormatXml

	// ReportFileExtExcel MicroSoft Excel file
	ReportFileExtExcel = ReportFileExtDot + "xlsx"
)

type ReportDataRequest struct {
	Id     string `json:"report_id,omitempty"`
	Format string `json:"output_format,omitempty"`
}

func ValidateReportDataRequest(req *ReportDataRequest) error {
	if req.Id == EmptyString {
		return errors.New("report id can not be empty")
	}
	switch req.Format {
	case ReportFormatJson, ReportFormatHtml, ReportFormatXml, ReportFormatExcel:
		return nil
	}
	return fmt.Errorf("unsupported report format: %s", req.Format)
}
