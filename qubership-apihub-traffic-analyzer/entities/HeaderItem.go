package entities

import (
	"crypto/md5"
	"fmt"
)

type HttpHeaderItem struct {
	tableName struct{} `pg:"http_headers, alias:http_headers"`

	Id    string `pg:"header_id, pk, type:varchar"`
	Key   string `pg:"name, type:varchar"`
	Value string `pg:"value, type:varchar"`
}

// ComputeHeaderId
// computes HTTP header checksum
func computeHeaderId(key, value string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key+value)))
}

// NewHttpHeader
// prepares new record to interact with DB
func NewHttpHeader(key, value string) HttpHeaderItem {
	return HttpHeaderItem{Key: key, Value: value, Id: computeHeaderId(key, value)}
}
