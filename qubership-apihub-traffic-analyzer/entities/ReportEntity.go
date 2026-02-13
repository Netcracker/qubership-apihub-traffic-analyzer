package entities

import (
	"time"
)

type ReportEntity struct {
	tableName struct{} `pg:"stored_reports, alias:stored_reports"`

	ReportId         int       `pg:"report_id,pk,type:BIGINT"`
	CreatedAt        time.Time `pg:"created_at,type:TIMESTAMP"`
	ReportParameters string    `pg:"report_parameters,type:json"`
	ReportTypeId     int       `pg:"report_type_id,type:INT"`
	ReportStatusId   int       `pg:"report_status_id,type:INT"`
	CompletedAt      time.Time `pg:"completed_at,type:TIMESTAMP"`
	ReportUuid       string    `pg:"report_uuid,type:varchar"`
}
