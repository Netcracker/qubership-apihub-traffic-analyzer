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
