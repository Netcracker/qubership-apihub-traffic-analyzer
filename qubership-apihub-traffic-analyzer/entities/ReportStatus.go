package entities

import (
	"time"
)

const (
	ReportStatusCreated = "created"
	ReportStatusFailed  = "failed"
	ReportStatusReady   = "ready"
)

type ReportStatusEntity struct {
	tableName struct{} `pg:"report_status, alias:report_status"`

	Id        int       `pg:"report_status_id,pk,type:SERIAL"`
	Name      string    `pg:"report_status,type:varchar,notnull"`
	CreatedAt time.Time `pg:"created_at,type:TIMESTAMP,notnull"`
	RetiredAt time.Time `pg:"retired_at,type:TIMESTAMP"`
}
