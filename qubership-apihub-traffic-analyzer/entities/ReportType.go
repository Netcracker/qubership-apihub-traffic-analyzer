package entities

import (
	"time"
)

const (
	ServiceOperationReport = "service operations"
)

type ReportTypeEntity struct {
	tableName struct{} `pg:"report_types, alias:report_types"`

	Id        int       `pg:"report_type_id,pk,type:INT"`
	Name      string    `pg:"report_type,type:varchar,notnull"`
	CreatedAt time.Time `pg:"created_at,type:TIMESTAMP,notnull"`
	RetiredAt time.Time `pg:"retired_at,type:TIMESTAMP"`
}
