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

import "encoding/json"

const (
	// OperationFound service operation found in capture
	OperationFound = "Captured API operation"
	// OperationNotFound service operation not found in capture
	OperationNotFound = "API operation not captured"
	// OperationExtra operation from dump does not belong to service
	OperationExtra = "Captured unknown request"
	// OperationDiff there are differences between service operation and dump (different method, path, ...)
	OperationDiff = "Different"
)

type OperationAddress struct {
	ServiceView
	IpAddress string `json:"ip_address,omitempty"`
	Port      int    `json:"tcp_port,omitempty"`
}
type OperationPeer struct {
	Source      OperationAddress `json:"source_address,omitempty"`
	Destination OperationAddress `json:"destination_address,omitempty"`
}

// OperationStatus
// an operation status for report
type OperationStatus struct {
	Id    string `json:"operation_id,omitempty"`
	Title string `json:"operation_title,omitempty"`
	Path  string `json:"operation_path,omitempty"`
	// Method operation call method (GET,POST,...)
	Method string `json:"operation_method,omitempty"`
	// Status an operation status (OperationFound, OperationNotFound, ...)
	Status   string          `json:"operation_status,omitempty"`
	HitCount int             `json:"dump_hit_count,omitempty"`
	Peers    []OperationPeer `json:"operation_peers,omitempty"`
}
type OperationStatusWithPeers struct {
	OperationStatus
	Source      string `json:"source_service,omitempty"`
	Destination string `json:"destination_service,omitempty"`
}

func DecodeOperationStatus(bytes []byte) (OperationStatus, error) {
	var status OperationStatus
	err := json.Unmarshal(bytes, &status)
	return status, err
}

func DecodeOperationStatusWithPeers(bytes []byte) (OperationStatusWithPeers, error) {
	var status OperationStatusWithPeers
	err := json.Unmarshal(bytes, &status)
	return status, err
}
