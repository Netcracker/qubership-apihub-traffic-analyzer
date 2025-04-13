// Copyright 2024-2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
