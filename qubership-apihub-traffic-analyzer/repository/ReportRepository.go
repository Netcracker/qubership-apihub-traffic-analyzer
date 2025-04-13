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
	"fmt"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/utils"
)

func SetReportParameters(report *entities.ReportEntity, reportParameters interface{}) error {
	jsonb, err := utils.MarshalToJSON(reportParameters)
	if err != nil {
		return err
	}
	report.ReportParameters = string(jsonb)
	return nil
}

func DeleteReport(db db.ConnectionProvider, report entities.ReportEntity) error {
	dataRows := new(entities.ReportDataRow)
	_, err := db.GetConnection().Model(dataRows).Where("report_id=?", report.ReportId).Delete()
	if err == nil {
		_, err = db.GetConnection().Model(&report).Delete()
	}
	return err
}

func UpdateReport(db db.ConnectionProvider, report entities.ReportEntity) error {
	_, err := db.GetConnection().Model(&report).WherePK().Update()
	return err
}

func InsertReport(db db.ConnectionProvider, report *entities.ReportEntity) error {
	result := new(entities.ReportEntity)
	_, err := db.GetConnection().Model(report).Returning("report_id, created_at").Insert(result)
	if err == nil {
		report.ReportId = result.ReportId
		report.CreatedAt = result.CreatedAt
	}
	return err
}

func GetReport(db db.ConnectionProvider, reportUuid, reportTypeName string) (*entities.ReportEntity, *entities.ReportTypeEntity, error) {
	result := new(entities.ReportEntity)
	err := db.GetConnection().Model(result).Where("report_uuid=?", reportUuid).First()
	if err == nil {
		reportStatus, err := GetReportStatusById(db, result.ReportStatusId)
		if err != nil {
			return result, nil, fmt.Errorf("unable to get report status for %d: %v", result.ReportStatusId, err)
		}
		if reportStatus.Name != entities.ReportStatusReady {
			return result, nil, fmt.Errorf("improper report status (%s instead of %s", reportStatus.Name, entities.ReportStatusReady)
		}
		reportType, err := GetReportTypeById(db, result.ReportTypeId)
		if err != nil {
			return result, reportType, fmt.Errorf("unable to get report type for %d: %v", result.ReportTypeId, err)
		}
		if reportType.Name != reportTypeName {
			return result, reportType, fmt.Errorf("improper report type %s instead of %s", reportStatus.Name, reportTypeName)
		}
		return result, reportType, nil
	}
	return result, nil, fmt.Errorf("unable to query report parameters: %v", err)
}
