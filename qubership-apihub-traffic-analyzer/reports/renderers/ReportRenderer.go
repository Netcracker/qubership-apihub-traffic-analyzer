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
	"fmt"
	"io"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/reports/generators"
)

type ReportRenderer interface {
	MakeReportHeader() error
	ProcessRows() error
	RenderRow(reportRow *entities.ReportDataRow) error
	MakeReportFooter() error
	FlushData(w io.Writer) error
	GetFileName() string
	Dispose()
}

func NewReportRenderer(db db.ConnectionProvider,
	req interface{},
	workDir string, reportTypeName generators.ReportType) (ReportRenderer, error) {
	switch reportTypeName {
	case generators.ServiceOperationReport:
		return NewServiceOperationsRenderer(db, req, workDir)
	}
	return nil, fmt.Errorf("report type: %s has not implemented yet", reportTypeName)
}
