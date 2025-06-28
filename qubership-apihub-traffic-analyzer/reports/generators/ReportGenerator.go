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

package generators

import (
	"fmt"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/client"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
)

type ReportGenerator interface {
	Generate(rq interface{}) error
}
type ReportType string

const (
	ServiceOperationReport ReportType = "service operations"
)

type ReportGeneratorParameters struct {
	ApihubClient  client.ApihubClient
	KubeNameSpace string
	WorkSpace     string
	WorkDir       string
	AgentName     string
	ReportType    ReportType
	Db            db.ConnectionProvider
}

func NewReportGenerator(parameters ReportGeneratorParameters) (ReportGenerator, error) {
	switch parameters.ReportType {
	case ServiceOperationReport:
		return NewServiceOperationsReport(parameters)
	}
	return nil, fmt.Errorf("report type: %s has not implemented yet", parameters.ReportType)
}
