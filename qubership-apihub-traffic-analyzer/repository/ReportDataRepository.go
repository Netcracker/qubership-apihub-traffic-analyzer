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

package repository

import (
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/utils"
)

func SetReportRowData(data *entities.ReportDataRow, rowData interface{}) error {
	jsonb, err := utils.MarshalToJSON(rowData)
	if err != nil {
		return err
	}
	data.ReportRow = string(jsonb)
	return nil
}

func InsertReportRow(db db.ConnectionProvider, data *entities.ReportDataRow) error {
	dataRows := new(entities.ReportDataRow)
	_, err := db.GetConnection().Model(data).Returning("report_row_id").Insert(dataRows)
	if err == nil {
		data.ReportRowId = dataRows.ReportRowId
	}
	return err
}
