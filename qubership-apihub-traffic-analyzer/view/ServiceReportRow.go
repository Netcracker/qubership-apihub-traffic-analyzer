package view

type ServiceReportRow struct {
	OperationName   string `json:"operation_name,omitempty"`
	OperationPath   string `json:"operation_path,omitempty"`
	DetectionStatus string `json:"detection_status,omitempty"`
}
