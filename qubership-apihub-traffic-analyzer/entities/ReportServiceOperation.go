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

package entities

import (
	"strings"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
)

const (
	ServiceOperationStar      = "*"
	ServiceOperationStarRegex = "[^\\/]+"
)

type ReportServiceOperation struct {
	tableName         struct{} `pg:"report_service_operations, alias:report_service_operations"`
	ReportOperationId int      `pg:"report_operation_id,pk,type:BIGINT"`
	ReportId          int      `pg:"report_id,type:BIGINT"`
	Title             string   `pg:"operation_title,type:varchar"`
	Path              string   `pg:"operation_path,type:varchar"`
	Regexp            string   `pg:"operation_path_re,type:varchar"`
	Method            string   `pg:"operation_method,type:varchar"`
	Status            string   `pg:"operation_status,type:varchar"`
	HitCount          int      `pg:"operation_hit_count,type:int"`
}

type ServiceOperationUpdate struct {
	tableName struct{} `pg:"ip_packets, alias:ip_packets"`
	PacketId  int      `pg:"packet_id,pk,type:bigint"`
	Path      string   `pg:"request_path,type:varchar"`
	Method    string   `pg:"request_method,type:varchar"`
	HitCount  int      `pg:"operation_hit_count,type:int"`
}

type ReportServiceOperationWithPeers struct {
	tableName   struct{} `pg:"report_service_operations2, alias:report_service_operations2"`
	ReportId    int      `pg:"report_id,type:bigint"`
	Sender      string   `pg:"src_peer, type:varchar"`
	Receiver    string   `pg:"dst_peer, type:varchar"`
	Path        string   `pg:"operation_path,type:varchar"`
	Method      string   `pg:"operation_method,type:varchar"`
	OperationId string   `pg:"operation_title, type:varchar"`
	Occurrences int      `pg:"hit_count,type:int"`
	Comment     string   `pg:"operation_status,type:varchar"`
}

func NewReportServiceOperation(reportId int, title, path, method, status string) ReportServiceOperation {
	ret := ReportServiceOperation{
		ReportId: reportId,
		Title:    title,
		Path:     path,
		Method:   strings.ToUpper(method),
		Status:   status,
		HitCount: 0,
	}
	pathRegex := strings.Replace(path, ServiceOperationStar, ServiceOperationStarRegex, -1)
	if pathRegex != path {
		ret.Regexp = "^" + pathRegex // from begin of the path
	}
	return ret
}

func InsertReportServiceOperation(db db.ConnectionProvider, data *ReportServiceOperation) error {
	res := new(ReportServiceOperation)
	_, err := db.GetConnection().Model(data).Returning("report_operation_id").Insert(res)
	if err == nil {
		data.ReportOperationId = res.ReportOperationId
	}
	return err
}
