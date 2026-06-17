package entities

type PacketHeader struct {
	tableName struct{} `pg:"service_packet_headers, alias:service_packet_headers"`

	HeaderId string `pg:"header_id, pk, type:varchar"`
	PacketId int    `pg:"packet_id, pk, type:bigint"`
}
