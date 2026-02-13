package view

import (
	"encoding/json"
)

type ServiceView struct {
	Name    string `json:"service_name,omitempty"`
	Version string `json:"service_version,omitempty"`
	Valid   bool
}

func UnmarshalServiceView(svcViewBytes []byte) (ServiceView, error) {
	svc := new(ServiceView)
	err := json.Unmarshal(svcViewBytes, svc)
	svc.Valid = (err == nil)
	return *svc, err
}
