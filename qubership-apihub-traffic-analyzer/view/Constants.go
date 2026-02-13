package view

const (
	EmptyString = ""
	// ApiKeyHeader - HTTP header name for API key
	ApiKeyHeader                = "api-key"
	MinioDeleteCapturePath      = "/api/v1/admin/capture/{captureId}/delete"
	LoadPath                    = "/api/v1/admin/capture/{captureId}/load"   // LoadPath - request data load/update capture data
	LoadStatusReportPath        = "/api/v1/admin/capture/{captureId}/status" // LoadStatusReportPath produce report, based on loaded data
	ServiceOperationsReportPath = "/api/v1/report/service/operations/generate"
	ServiceOperationsRenderPath = "/api/v1/report/service/operations/render"
	MinioCleanupCapturePath     = "/api/v1/admin/capture/S3/cleanup"
	CaptureIdParam              = "captureId"
	CompressedSuffix            = ".gz"
	AddressListSuffix           = "_address_list.txt"
	CaptureSuffix               = ".pcap"
	MetadataSuffix              = "_metadata.json"
)
