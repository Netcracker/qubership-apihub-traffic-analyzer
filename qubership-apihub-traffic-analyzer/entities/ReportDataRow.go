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

type ReportDataRow struct {
	tableName   struct{} `pg:"report_data,alias:report_data"`
	ReportRowId int      `pg:"report_row_id,pk,type:BIGSERIAL"`
	ReportId    int      `pg:"report_id,type:BIGINT"`
	ReportRow   string   `pg:"report_row,type:json"`
}
