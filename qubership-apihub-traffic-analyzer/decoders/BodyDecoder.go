package decoders

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
)

type BodyReadResult struct {
	// body as byte array
	Body []byte
	// decoded length
	Length int
	// to store body decoding error
	Err error
}

// BodyToString
// decodes HTTP request or HTTP response body from reader
func BodyToString(Body io.Reader, Uncompressed bool) BodyReadResult {
	res := BodyReadResult{
		Body:   nil,
		Length: -1,
		Err:    nil,
	}
	body, err := io.ReadAll(Body)
	if err != nil {
		res.Err = fmt.Errorf("error reading body: %v", err)
	}
	if !Uncompressed {
		zr, err1 := gzip.NewReader(bytes.NewReader(body))
		if err1 == nil {
			res.Body, err = io.ReadAll(zr)
			_ = zr.Close()
			if err == nil {
				res.Length = len(res.Body)
			} else {
				res.Err = fmt.Errorf("unable to read all bytes: %w", err)
			}
		} else {
			res.Err = fmt.Errorf("unable to create gzip reader: %w", err1)
		}
	}
	if res.Body == nil && body != nil {
		res.Body = body
		res.Length = len(body)
	}
	if errors.Is(res.Err, io.EOF) {
		res.Err = nil
	}
	return res
}
