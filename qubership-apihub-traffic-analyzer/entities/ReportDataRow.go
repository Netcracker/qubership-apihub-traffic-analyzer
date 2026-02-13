package entities

type ReportDataRow struct {
	tableName   struct{} `pg:"report_data,alias:report_data"`
	ReportRowId int      `pg:"report_row_id,pk,type:BIGSERIAL"`
	ReportId    int      `pg:"report_id,type:BIGINT"`
	ReportRow   string   `pg:"report_row,type:json"`
}
