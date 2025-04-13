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

const OpenAPI31Type string = "openapi-3-1"
const OpenAPI30Type string = "openapi-3-0"
const OpenAPI20Type string = "openapi-2-0"
const AsyncAPIType string = "asyncapi-2"
const JsonSchemaType string = "json-schema"
const MDType string = "markdown"
const GraphQLSchemaType string = "graphql-schema"
const GraphAPIType string = "graphapi"
const GraphQLType string = "graphql"
const IntrospectionType string = "introspection"
const UnknownType string = "unknown"

type Specification struct {
	Name     string `json:"name"`
	Path     string `json:"-"`
	Format   string `json:"format"` // json or yaml
	FileId   string `json:"fileId"`
	Type     string `json:"type"`
	XApiKind string `json:"xApiKind,omitempty"`
}
