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
	"time"
)

type LoadHistoryKey struct {
	CaptureId string `json:"capture_id"`
}

type LoadHistoryValue struct {
	BeginDateTime string `json:"begin_date_time"`
	EndDateTime   string `json:"end_date_time"`
	Error         error  `json:"error"`
}

type LoadHistoryRecord struct {
	LoadHistoryKey
	LoadHistoryValue
}

// GetHistoryDateTimeString
// returns date and time string to fill BeginDateTime/EndDateTime
func GetHistoryDateTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// HistoryRecordCompleted
// returns record completeness status and error result
func HistoryRecordCompleted(val LoadHistoryValue) (bool, error) {
	return val.EndDateTime != EmptyString, val.Error
}
