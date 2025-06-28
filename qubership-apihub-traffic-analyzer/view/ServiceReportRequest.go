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
	"encoding/json"
	"errors"
)

type ServiceReportRequest struct {
	// a unique report_id received at report creation
	ReportUuid string `json:"report_uuid,omitempty"`
	// a capture id used to create the report
	CaptureId string `json:"capture_id"`
	// a service name to validate operation from
	ServiceName string `json:"service_name"`
	// a service version used to receive operation list
	ServiceVersion string `json:"service_version,omitempty"`
	// a version status (requested or the most recent at APIHUB)
	VersionStatus string `json:"version_status,omitempty"`
}

func ValidateServiceReportRequest(req ServiceReportRequest) error {
	if req.CaptureId == EmptyString {
		return errors.New("capture_id is empty")
	}
	if req.ServiceName == EmptyString {
		return errors.New("service_name is empty")
	}
	return nil
}

func UnmarshalServiceReportRequest(svcViewBytes []byte) (ServiceReportRequest, error) {
	svc := new(ServiceReportRequest)
	err := json.Unmarshal(svcViewBytes, svc)
	return *svc, err
}
