package entities

import (
	"encoding/json"
)

type ServiceAddress struct {
	tableName struct{} `pg:"service_addresses, alias:service_addresses"`

	Id        int    `pg:"address_id, pk, type:bigint"`
	Address   string `pg:"ip_address, type:varchar"`
	Name      string `pg:"service_name, type:varchar"`
	Version   string `pg:"service_version, type:varchar"`
	CaptureId string `pg:"capture_id, type:varchar"`
}

func UnmarshalServiceAddress(bytes []byte) (*ServiceAddress, error) {
	result := new(ServiceAddress)
	err := json.Unmarshal(bytes, result)
	return result, err
}
