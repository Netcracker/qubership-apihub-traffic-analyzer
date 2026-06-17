package entities

const (
	ReportAffectedPacket    = 1
	ReportAffectedOperation = 2
)

type ReportAffectedRef struct {
	tableName     struct{} `pg:"report_affected_rows, alias:report_affected_rows"`
	ReportId      int      `pg:"report_id,type:BIGINT"`
	ReferenceId   int      `pg:"reference_id,type:BIGINT"`
	ReferenceType int      `pg:"reference_type,type:int"`
	HitCount      int      `pg:"hit_count,type:int"`
}
