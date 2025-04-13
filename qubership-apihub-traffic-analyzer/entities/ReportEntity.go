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
	"time"
)

type ReportEntity struct {
	tableName struct{} `pg:"stored_reports, alias:stored_reports"`

	ReportId         int       `pg:"report_id,pk,type:BIGINT"`
	CreatedAt        time.Time `pg:"created_at,type:TIMESTAMP"`
	ReportParameters string    `pg:"report_parameters,type:json"`
	ReportTypeId     int       `pg:"report_type_id,type:INT"`
	ReportStatusId   int       `pg:"report_status_id,type:INT"`
	CompletedAt      time.Time `pg:"completed_at,type:TIMESTAMP"`
	ReportUuid       string    `pg:"report_uuid,type:varchar"`
}
