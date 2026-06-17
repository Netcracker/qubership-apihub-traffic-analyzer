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
