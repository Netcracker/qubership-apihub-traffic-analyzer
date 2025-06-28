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

import (
	"fmt"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/openapi/orderedmap"
)

type Operation struct {
	OperationId             string                 `json:"operationId" validate:"required"`
	Title                   string                 `json:"title" validate:"required"`
	ApiType                 string                 `json:"apiType" validate:"required"`
	DataHash                string                 `json:"dataHash" validate:"required"`
	Deprecated              bool                   `json:"deprecated"`
	ApiKind                 string                 `json:"apiKind" validate:"required"`
	Metadata                map[string]interface{} `json:"metadata" validate:"required"`
	SearchScopes            map[string]interface{} `json:"searchScopes" validate:"required"`
	PreviousReleaseVersions []string               `json:"deprecatedInPreviousVersions"`
	DeprecatedInfo          map[string]interface{} `json:"deprecatedInfo"`
	DeprecatedItems         []DeprecatedItem       `json:"deprecatedItems"`
	Tags                    []string               `json:"tags"`
	Models                  map[string]string      `json:"models"`
}

type DocumentsOperation struct {
	OperationId string      `json:"operationId" validate:"required"`
	Title       string      `json:"title" validate:"required"`
	ApiType     string      `json:"apiType" validate:"required"`
	DataHash    string      `json:"dataHash" validate:"required"`
	Deprecated  bool        `json:"deprecated"`
	ApiKind     string      `json:"apiKind" validate:"required"`
	Metadata    interface{} `json:"metadata" validate:"required"`
}

type SingleOperationView struct {
	Data        *orderedmap.OrderedMap `json:"data,omitempty"`
	OperationId string                 `json:"operationId"`
	Title       string                 `json:"title"`
	DataHash    string                 `json:"dataHash"`
	Deprecated  bool                   `json:"deprecated,omitempty"`
	ApiKind     string                 `json:"apiKind"`
	ApiType     string                 `json:"apiType"`
	CustomTags  map[string]interface{} `json:"customTags,omitempty"`
}

type CommonOperationView struct {
	OperationId string                 `json:"operationId"`
	Title       string                 `json:"title"`
	DataHash    string                 `json:"dataHash"`
	Deprecated  bool                   `json:"deprecated,omitempty"`
	ApiKind     string                 `json:"apiKind"`
	ApiType     string                 `json:"apiType"`
	CustomTags  map[string]interface{} `json:"customTags,omitempty"`
}

type OperationListView struct {
	CommonOperationView
	PackageRef string                 `json:"packageRef,omitempty"`
	Data       *orderedmap.OrderedMap `json:"data,omitempty"`
}

type DeprecatedOperationView struct {
	PackageRef              string                 `json:"packageRef,omitempty"`
	OperationId             string                 `json:"operationId"`
	Title                   string                 `json:"title"`
	DataHash                string                 `json:"dataHash"`
	Deprecated              bool                   `json:"deprecated,omitempty"`
	ApiKind                 string                 `json:"apiKind"`
	ApiType                 string                 `json:"apiType"`
	PreviousReleaseVersions []string               `json:"deprecatedInPreviousVersions,omitempty"`
	DeprecatedCount         int                    `json:"deprecatedCount"`
	DeprecatedInfo          map[string]interface{} `json:"deprecatedInfo,omitempty"`
	DeprecatedItems         []DeprecatedItem       `json:"deprecatedItems,omitempty"`
}
type DeprecatedItem struct {
	PreviousReleaseVersions []string               `json:"deprecatedInPreviousVersions,omitempty"`
	JsonPath                []string               `json:"jsonPath,omitempty"`
	Description             string                 `json:"description,omitempty"`
	DeprecatedInfo          map[string]interface{} `json:"deprecatedInfo,omitempty"`
}

type DeprecatedItems struct {
	DeprecatedItems []DeprecatedItem `json:"deprecatedItems"`
}
type OperationChangelogView struct {
	OperationId string                  `json:"operationId"`
	Title       string                  `json:"title"`
	Deprecated  bool                    `json:"deprecated,omitempty"`
	ApiKind     string                  `json:"apiKind"`
	Action      string                  `json:"action"`
	Changes     []SingleOperationChange `json:"changes"`
}

type OperationComparison struct {
	OperationId      string                  `json:"operationId" validate:"required"`
	DataHash         string                  `json:"dataHash,omitempty"`
	PreviousDataHash string                  `json:"previousDataHash,omitempty"`
	ChangeSummary    ChangeSummary           `json:"changeSummary,omitempty"`
	Changes          []SingleOperationChange `json:"changes" validate:"required,dive,required"`
	JsonPath         []string                `json:"jsonPath,omitempty"`
	Action           string                  `json:"action,omitempty"`
	Severity         string                  `json:"severity,omitempty"`
	Metadata         map[string]interface{}  `json:"metadata"`
}
type SingleOperationChange struct {
	Path        []string `json:"jsonPath,omitempty"`
	Action      string   `json:"action,omitempty"`
	Severity    string   `json:"severity,omitempty"`
	Description string   `json:"description,omitempty"`
}

type VersionChangesView struct {
	PreviousVersion          string                       `json:"previousVersion"`
	PreviousVersionPackageId string                       `json:"previousVersionPackageId"`
	Operations               []interface{}                `json:"operations"`
	Packages                 map[string]PackageVersionRef `json:"packages,omitempty"`
}

type OperationComparisonChangelogView_deprecated struct {
	OperationId               string        `json:"operationId"`
	Title                     string        `json:"title"`
	ApiKind                   string        `json:"apiKind,omitempty"`
	DataHash                  string        `json:"dataHash,omitempty"`
	PreviousDataHash          string        `json:"previousDataHash,omitempty"`
	ChangeSummary             ChangeSummary `json:"changeSummary"`
	PackageRef                string        `json:"packageRef"`
	PreviousVersionPackageRef string        `json:"previousVersionPackageRef"`
}

type ComparisonOperationView struct {
	Title      string `json:"title"`
	ApiKind    string `json:"apiKind,omitempty"`
	DataHash   string `json:"dataHash,omitempty"`
	PackageRef string `json:"packageRef"`
}

type OperationComparisonChangelogView struct {
	OperationId       string                   `json:"operationId"`
	CurrentOperation  *ComparisonOperationView `json:"currentOperation,omitempty"`
	PreviousOperation *ComparisonOperationView `json:"previousOperation,omitempty"`
	ChangeSummary     ChangeSummary            `json:"changeSummary"`
}

type OperationComparisonChangesView struct {
	OperationId               string                  `json:"operationId"`
	Title                     string                  `json:"title"`
	ApiKind                   string                  `json:"apiKind,omitempty"`
	DataHash                  string                  `json:"dataHash,omitempty"`
	PreviousDataHash          string                  `json:"previousDataHash,omitempty"`
	ChangeSummary             ChangeSummary           `json:"changeSummary"`
	PackageRef                string                  `json:"packageRef"`
	PreviousVersionPackageRef string                  `json:"previousVersionPackageRef"`
	Changes                   []SingleOperationChange `json:"changes"`
	Action                    string                  `json:"action"`
}

type OperationChangesView struct {
	Changes []SingleOperationChange `json:"changes"`
}

type OperationTags struct {
	Tags []string `json:"tags"`
}

type Operations struct {
	Operations []interface{}                `json:"operations"`
	Packages   map[string]PackageVersionRef `json:"packages,omitempty"`
}

type GroupedOperations struct {
	Operations []interface{}                `json:"operations"`
	Packages   map[string]PackageVersionRef `json:"packages,omitempty"`
}

type GroupedGhostOperations_deprecated struct {
	GhostOperations []interface{}                `json:"ghostOperations"`
	Packages        map[string]PackageVersionRef `json:"packages,omitempty"`
}

type ChangeSummary struct {
	Breaking     int `json:"breaking"`
	SemiBreaking int `json:"semi-breaking"`
	Deprecated   int `json:"deprecated"`
	NonBreaking  int `json:"non-breaking"`
	Annotation   int `json:"annotation"`
	Unclassified int `json:"unclassified"`
}

func (c ChangeSummary) GetTotalSummary() int {
	return c.Breaking + c.SemiBreaking + c.Deprecated + c.NonBreaking + c.Annotation + c.Unclassified
}

const ChangelogActionChange string = "change"
const ChangelogActionAdd string = "add"
const ChangelogActionRemove string = "remove"

type ApiKind string

const BwcApiKind ApiKind = "bwc"
const NoBwcApiKind ApiKind = "no-bwc"
const DebugApiKind ApiKind = "debug"
const ExperimentalApiKind ApiKind = "experimental"

type Severity string

const Annotation Severity = "annotation"
const Breaking Severity = "breaking"
const SemiBreaking Severity = "semi-breaking"
const Deprecated Severity = "deprecated"
const NonBreaking Severity = "non-breaking"
const Unclassified Severity = "unclassified"

func ValidSeverity(s string) bool {
	switch s {
	case string(Annotation), string(Breaking), string(SemiBreaking), string(Deprecated), string(NonBreaking), string(Unclassified):
		return true
	}
	return false
}

func ParseApiKind(s string) (ApiKind, error) {
	switch s {
	case string(BwcApiKind):
		return BwcApiKind, nil
	case string(NoBwcApiKind):
		return NoBwcApiKind, nil
	case string(DebugApiKind):
		return DebugApiKind, nil
	case string(ExperimentalApiKind):
		return ExperimentalApiKind, nil
	default:
		return "", fmt.Errorf("unknown API Kind: %v", s)
	}
}

type ApiType string

const RestApiType ApiType = "rest"
const GraphqlApiType ApiType = "graphql"

func ParseApiType(s string) (ApiType, error) {
	switch s {
	case string(RestApiType):
		return RestApiType, nil
	case string(GraphqlApiType):
		return GraphqlApiType, nil
	default:
		return "", fmt.Errorf("unknown API Type: %v", s)
	}
}

func GetDocumentTypesForApiType(apiType string) []string {
	switch apiType {
	case string(RestApiType):
		return []string{OpenAPI20Type, OpenAPI30Type, OpenAPI31Type}
	case string(GraphqlApiType):
		return []string{GraphQLSchemaType, GraphAPIType, IntrospectionType}
	default:
		return []string{}
	}
}

type OperationListReq struct {
	Deprecated     *bool
	HashList       []string
	Ids            []string
	IncludeData    bool
	Kind           string
	EmptyTag       bool
	Tag            string
	Limit          int
	Page           int
	TextFilter     string
	ApiType        string
	DocumentSlug   string
	EmptyGroup     bool
	Group          string
	OnlyAddable    bool
	RefPackageId   string
	CustomTagKey   string
	CustomTagValue string
}

type DeprecatedOperationListReq struct {
	Ids                    []string
	Kind                   string
	Tags                   []string
	Limit                  int
	Page                   int
	TextFilter             string
	ApiType                string
	DocumentSlug           string
	IncludeDeprecatedItems bool
	RefPackageId           string
	EmptyTag               bool
	EmptyGroup             bool
	Group                  string
}

type OperationBasicSearchReq struct {
	PackageId   string
	Version     string
	OperationId string
	Revision    int
	ApiType     string
	ApiKind     string
	Limit       int
	Offset      int
	TextFilter  string
}

type VersionChangesReq struct {
	PreviousVersion          string
	PreviousVersionPackageId string
	DocumentSlug             string
	ApiKind                  string
	EmptyTag                 bool
	RefPackageId             string
	Tags                     []string
	TextFilter               string
	Limit                    int
	Offset                   int
	EmptyGroup               bool
	Group                    string
	Severities               []string
}

type PagingFilterReq struct {
	Limit      int
	Offset     int
	TextFilter string
}

type DocumentsFilterReq struct {
	ApiType    string
	Limit      int
	Offset     int
	TextFilter string
}

type DocumentsForTransformationFilterReq struct {
	ApiType                string
	Limit                  int
	Offset                 int
	FilterByOperationGroup string
}
type ChangelogCalculationParams struct {
	PreviousVersion          string `json:"previousVersion"`
	PreviousVersionPackageId string `json:"previousVersionPackageId"`
	ReCalculate              bool   `json:"reCalculate"`
}

type CalculationProcessStatus struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message"`
}

type DeprecatedOperationsSummary struct {
	OperationTypes *[]DeprecatedOperationType     `json:"operationTypes,omitempty"`
	Refs           *[]DeprecatedOperationTypesRef `json:"refs,omitempty"`
	Packages       *map[string]PackageVersionRef  `json:"packages,omitempty"`
}
type DeprecatedOperationType struct {
	ApiType         string   `json:"apiType"`
	DeprecatedCount int      `json:"deprecatedCount"`
	Tags            []string `json:"tags"`
}
type DeprecatedOperationTypesRef struct {
	PackageRef     string                    `json:"packageRef,omitempty"`
	OperationTypes []DeprecatedOperationType `json:"operationTypes"`
}

type OperationModelUsages struct {
	ModelUsages []OperationModels `json:"modelUsages"`
}

type OperationModels struct {
	OperationId string   `json:"operationId"`
	ModelNames  []string `json:"modelNames"`
}

type CustomTags struct {
	CustomTag map[string]interface{}
}
