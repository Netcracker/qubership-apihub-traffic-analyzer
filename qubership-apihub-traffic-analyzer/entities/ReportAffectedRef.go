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

const (
	ReportAffectedPacket    = 1
	ReportAffectedOperation = 2
)

type ReportAffectedRef struct {
	tableName     struct{} `pg:"report_affected_rows, alias:report_affected_rows"`
	ReportId      int      `pg:"report_id,type:BIGINT"`
	ReferenceId   int      `pg:"reference_id,type:BIGINT"`
	ReferenceType int      `pg:"reference_type,type:int"`
	HitCount      int      `pg:"hit_count,type:int"`
}
