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
