package view

import (
	"strconv"
)

const (
	SourcePeer = 0
	DestPeer   = 1
)

type EndPointView struct {
	Address string `json:"ip_address,omitempty"`
	Port    int    `json:"port,omitempty"`
	ServiceView
}

func EndPointToString(endPoint EndPointView) string {
	return endPoint.Address + ":" + strconv.Itoa(endPoint.Port)
}
