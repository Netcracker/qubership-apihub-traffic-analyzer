package entities

type CaptureMetadata struct {
	tableName struct{} `pg:"capture_metadata, alias:capture_metadata"`

	CaptureId string `pg:"capture_id,pk,type:varchar"`
	Metadata  string `pg:"capture_metadata,type:json"`
}
