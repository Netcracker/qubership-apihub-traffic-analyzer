package view

import (
	"time"
)

type LoadHistoryKey struct {
	CaptureId string `json:"capture_id"`
}

type LoadHistoryValue struct {
	BeginDateTime string `json:"begin_date_time"`
	EndDateTime   string `json:"end_date_time"`
	Error         error  `json:"error"`
}

type LoadHistoryRecord struct {
	LoadHistoryKey
	LoadHistoryValue
}

// GetHistoryDateTimeString
// returns date and time string to fill BeginDateTime/EndDateTime
func GetHistoryDateTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// HistoryRecordCompleted
// returns record completeness status and error result
func HistoryRecordCompleted(val LoadHistoryValue) (bool, error) {
	return val.EndDateTime != EmptyString, val.Error
}
