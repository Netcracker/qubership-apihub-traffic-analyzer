package decoders

import (
	"regexp"
)

type PayloadType int

const (
	// PTNotHttp - not a HTTP packet
	PTNotHttp PayloadType = iota
	// PTHttpReq - expected HTTP request
	PTHttpReq
	// PTHttpResp - expected HTTP response
	PTHttpResp
	// PTHttp supposed to be generic HTTP packet
	PTHttp
	byteH              = byte('H')
	byteP              = byte('P')
	byteT              = byte('T')
	RequestMethod      = "Method"
	RequestPath        = "Path"
	RRProtocol         = "Proto"
	ResponseStatus     = "Status"
	ResponseStatusText = "StatusText"
	// RRMatchNum an index of the match
	RRMatchNum = "MatchPosition"
	// RHOffset an offset of the HTTP marker
	RHOffset = "Offset"
	/*RRPayloadType      = "PayloadType"*/
)

// httpB
// a Hypertext protocol marker fallback
var httpB = [4]byte{byteH, byteT, byteT, byteP}

// DecodeFeedback
// a decoder feedback interface implementation to use libpcap frame decoders
type DecodeFeedback struct {
}

// SetTruncated
// to make decoder feedback interface implementation valid
func (df *DecodeFeedback) SetTruncated() {

}

var (
	// reqRe HTTP request decoder regular expression
	reqRe = regexp.MustCompile(`^\s*(\w+)\s+(\S+)\s+(HTTP\/\d+\.\d+)\s`)
	// respRe HTTP response decoder regular expression
	respRe = regexp.MustCompile(`(?m)^\s*(HTTP\/\d+\.\d+)\s+(\d+)\s+([\w\s_.]+)\s+`)
)

// DetectHttp
// tries to detect HTTP data in the given byte array, returns detection result and data (if detected)
func DetectHttp(payLoad []byte) (PayloadType, map[string]interface{}) {
	fields := make(map[string]interface{})
	// is it an HTTP request?
	for i, matches := range reqRe.FindAllSubmatch(payLoad, 1) {
		if len(matches) > 1 {
			for j, m := range matches {
				switch j {
				case 1:
					fields[RequestMethod] = string(m)
				case 2:
					fields[RequestPath] = string(m)
				case 3:
					fields[RRProtocol] = string(m)
				default:
					break
				}
			}
			fields[RRMatchNum] = i
			return PTHttpReq, fields // that's a request
		}

	}
	// is it an HTTP response?
	for i, matches := range respRe.FindAllSubmatch(payLoad, 1) {
		if len(matches) > 1 {
			for j, m := range matches {
				switch j {
				case 1:
					fields[RRProtocol] = string(m)
				case 2:
					fields[ResponseStatus] = string(m)
				case 3:
					fields[ResponseStatusText] = string(m)
				default:
					break
				}
			}
			fields[RRMatchNum] = i
			return PTHttpResp, fields // that's a response
		}

	}
	// not a request/response - try to find anything like HTTP
	pos := 0
	pl := len(payLoad)
	hl := len(httpB) - 1
	bFound := false
	for pos < pl {
		for i := 0; i <= hl && pos < pl; i++ {
			if payLoad[pos] == httpB[i] {
				pos++
				if i == hl {
					bFound = true
				}
			} else {
				pos++
				break
			}
		}
		if bFound {
			break
		}
	}
	if bFound {
		if pos >= pl {
			pos = -1
		}
		fields[RHOffset] = pos
		return PTHttp, fields // something an HTTP-like
	}
	return PTNotHttp, fields // this is not an HTTP data
}
