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

import "time"

type PackagesInfo struct {
	Id                        string              `json:"packageId"`
	Alias                     string              `json:"alias"`
	ParentId                  string              `json:"parentId"`
	Kind                      string              `json:"kind"`
	Name                      string              `json:"name"`
	Description               string              `json:"description"`
	IsFavorite                bool                `json:"isFavorite"`
	ServiceName               string              `json:"serviceName,omitempty"`
	ImageUrl                  string              `json:"imageUrl"`
	Parents                   []ParentPackageInfo `json:"parents"`
	DefaultRole               string              `json:"defaultRole"`
	UserPermissions           []string            `json:"permissions"`
	LastReleaseVersionDetails *VersionDetails     `json:"lastReleaseVersionDetails,omitempty"`
	RestGroupingPrefix        string              `json:"restGroupingPrefix,omitempty"`
}

type Packages struct {
	Packages []PackagesInfo `json:"packages"`
}

type VersionDetails struct {
	Version           string         `json:"version"`
	NotLatestRevision bool           `json:"notLatestRevision,omitempty"`
	Summary           *ChangeSummary `json:"summary,omitempty"`
}

type PackagesSearchReq struct {
	ServiceName        string
	TextFilter         string
	Kind               string
	ParentId           string
	ShowAllDescendants bool
	Page               int
	Limit              int
}

type SimplePackage struct {
	Id                    string              `json:"packageId"`
	Alias                 string              `json:"alias"`
	ParentId              string              `json:"parentId"`
	Kind                  string              `json:"kind"`
	Name                  string              `json:"name"`
	Description           string              `json:"description"`
	IsFavorite            bool                `json:"isFavorite"`
	ServiceName           string              `json:"serviceName,omitempty"`
	ImageUrl              string              `json:"imageUrl"`
	Parents               []ParentPackageInfo `json:"parents"`
	UserRole              string              `json:"userRole"`
	DefaultRole           string              `json:"defaultRole"`
	DeletionDate          *time.Time          `json:"-"`
	DeletedBy             string              `json:"-"`
	CreatedBy             string              `json:"-"`
	CreatedAt             time.Time           `json:"-"`
	ReleaseVersionPattern string              `json:"releaseVersionPattern"`
}

type SimplePackages struct {
	Packages []SimplePackage `json:"packages"`
}
