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
