package view

import "fmt"

type ServiceStruct struct {
	Id                       string            `json:"id"`
	Name                     string            `json:"serviceName"`
	Url                      string            `json:"url"`
	Specs                    []Specification   `json:"specs"`
	Baseline                 *Baseline         `json:"baseline,omitempty"`
	Labels                   map[string]string `json:"serviceLabels,omitempty"`
	AvailablePromoteStatuses []string          `json:"availablePromoteStatuses"`
	ProxyServerUrl           string            `json:"proxyServerUrl,omitempty"`
}

type StatusEnum string

const StatusNone StatusEnum = "none"
const StatusRunning StatusEnum = "running"
const StatusComplete StatusEnum = "complete"
const StatusError StatusEnum = "error"
const StatusFailed StatusEnum = "failed"

type ServiceListResponse struct {
	Services []ServiceStruct `json:"services"`
	Status   StatusEnum      `json:"status"`
	Debug    string          `json:"debug"`
}

type ServiceNameItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ServiceNamesResponse struct {
	ServiceNames []ServiceNameItem `json:"serviceNames"`
}

type Baseline struct {
	PackageId string   `json:"packageId"`
	Name      string   `json:"name"`
	Url       string   `json:"url"`
	Versions  []string `json:"versions"`
}

func BuildStatusFromString(str string) (StatusEnum, error) {
	switch str {
	case "none":
		return StatusNone, nil
	case "running":
		return StatusRunning, nil
	case "complete":
		return StatusComplete, nil
	case "error":
		return StatusError, nil
	case "failure":
		return StatusFailed, nil
	}
	return StatusNone, fmt.Errorf("unknown build status: %s", str)
}
