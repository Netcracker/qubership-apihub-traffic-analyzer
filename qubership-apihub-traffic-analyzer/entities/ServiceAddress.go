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
