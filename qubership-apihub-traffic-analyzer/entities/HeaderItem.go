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
	"crypto/md5"
	"fmt"
)

type HttpHeaderItem struct {
	tableName struct{} `pg:"http_headers, alias:http_headers"`

	Id    string `pg:"header_id, pk, type:varchar"`
	Key   string `pg:"name, type:varchar"`
	Value string `pg:"value, type:varchar"`
}

// ComputeHeaderId
// computes HTTP header checksum
func computeHeaderId(key, value string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key+value)))
}

// NewHttpHeader
// prepares new record to interact with DB
func NewHttpHeader(key, value string) HttpHeaderItem {
	return HttpHeaderItem{Key: key, Value: value, Id: computeHeaderId(key, value)}
}
