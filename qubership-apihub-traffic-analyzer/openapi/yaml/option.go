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

package yaml

import (
	"io"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/openapi/yaml/ast"
)

// DecodeOption functional option type for Decoder
type DecodeOption func(d *Decoder) error

// ReferenceReaders pass to Decoder that reference to anchor defined by passed readers
func ReferenceReaders(readers ...io.Reader) DecodeOption {
	return func(d *Decoder) error {
		d.referenceReaders = append(d.referenceReaders, readers...)
		return nil
	}
}

// ReferenceFiles pass to Decoder that reference to anchor defined by passed files
func ReferenceFiles(files ...string) DecodeOption {
	return func(d *Decoder) error {
		d.referenceFiles = files
		return nil
	}
}

// ReferenceDirs pass to Decoder that reference to anchor defined by files under the passed dirs
func ReferenceDirs(dirs ...string) DecodeOption {
	return func(d *Decoder) error {
		d.referenceDirs = dirs
		return nil
	}
}

// RecursiveDir search yaml file recursively from passed dirs by ReferenceDirs option
func RecursiveDir(isRecursive bool) DecodeOption {
	return func(d *Decoder) error {
		d.isRecursiveDir = isRecursive
		return nil
	}
}

// Validator set StructValidator instance to Decoder
func Validator(v StructValidator) DecodeOption {
	return func(d *Decoder) error {
		d.validator = v
		return nil
	}
}

// Strict enable DisallowUnknownField and DisallowDuplicateKey
func Strict() DecodeOption {
	return func(d *Decoder) error {
		d.disallowUnknownField = true
		d.disallowDuplicateKey = true
		return nil
	}
}

// DisallowUnknownField causes the Decoder to return an error when the destination
// is a struct and the input contains object keys which do not match any
// non-ignored, exported fields in the destination.
func DisallowUnknownField() DecodeOption {
	return func(d *Decoder) error {
		d.disallowUnknownField = true
		return nil
	}
}

// DisallowDuplicateKey causes an error when mapping keys that are duplicates
func DisallowDuplicateKey() DecodeOption {
	return func(d *Decoder) error {
		d.disallowDuplicateKey = true
		return nil
	}
}

// UseOrderedMap can be interpreted as a map,
// and uses MapSlice ( ordered map ) aggressively if there is no type specification
func UseOrderedMap() DecodeOption {
	return func(d *Decoder) error {
		d.useOrderedMap = true
		return nil
	}
}

// UseJSONUnmarshaler if neither `BytesUnmarshaler` nor `InterfaceUnmarshaler` is implemented
// and `UnmashalJSON([]byte)error` is implemented, convert the argument from `YAML` to `JSON` and then call it.
func UseJSONUnmarshaler() DecodeOption {
	return func(d *Decoder) error {
		d.useJSONUnmarshaler = true
		return nil
	}
}

// EncodeOption functional option type for Encoder
type EncodeOption func(e *Encoder) error

// Indent change indent number
func Indent(spaces int) EncodeOption {
	return func(e *Encoder) error {
		e.indent = spaces
		return nil
	}
}

// IndentSequence causes sequence values to be indented the same value as Indent
func IndentSequence(indent bool) EncodeOption {
	return func(e *Encoder) error {
		e.indentSequence = indent
		return nil
	}
}

// UseSingleQuote determines if single or double quotes should be preferred for strings.
func UseSingleQuote(sq bool) EncodeOption {
	return func(e *Encoder) error {
		e.singleQuote = sq
		return nil
	}
}

// Flow encoding by flow style
func Flow(isFlowStyle bool) EncodeOption {
	return func(e *Encoder) error {
		e.isFlowStyle = isFlowStyle
		return nil
	}
}

// UseLiteralStyleIfMultiline causes encoding multiline strings with a literal syntax,
// no matter what characters they include
func UseLiteralStyleIfMultiline(useLiteralStyleIfMultiline bool) EncodeOption {
	return func(e *Encoder) error {
		e.useLiteralStyleIfMultiline = useLiteralStyleIfMultiline
		return nil
	}
}

// JSON encode in JSON format
func JSON() EncodeOption {
	return func(e *Encoder) error {
		e.isJSONStyle = true
		e.isFlowStyle = true
		return nil
	}
}

// MarshalAnchor call back if encoder find an anchor during encoding
func MarshalAnchor(callback func(*ast.AnchorNode, interface{}) error) EncodeOption {
	return func(e *Encoder) error {
		e.anchorCallback = callback
		return nil
	}
}

// UseJSONMarshaler if neither `BytesMarshaler` nor `InterfaceMarshaler`
// nor `encoding.TextMarshaler` is implemented and `MarshalJSON()([]byte, error)` is implemented,
// call `MarshalJSON` to convert the returned `JSON` to `YAML` for processing.
func UseJSONMarshaler() EncodeOption {
	return func(e *Encoder) error {
		e.useJSONMarshaller = true
		return nil
	}
}

// CommentPosition type of the position for comment.
type CommentPosition int

const (
	CommentLinePosition CommentPosition = iota
	CommentHeadPosition
)

func (p CommentPosition) String() string {
	switch p {
	case CommentLinePosition:
		return "Line"
	case CommentHeadPosition:
		return "Head"
	default:
		return ""
	}
}

// LineComment create a one-line comment for CommentMap.
func LineComment(text string) *Comment {
	return &Comment{
		Texts:    []string{text},
		Position: CommentLinePosition,
	}
}

// HeadComment create a multiline comment for CommentMap.
func HeadComment(texts ...string) *Comment {
	return &Comment{
		Texts:    texts,
		Position: CommentHeadPosition,
	}
}

// Comment raw data for comment.
type Comment struct {
	Texts    []string
	Position CommentPosition
}

// CommentMap map of the position of the comment and the comment information.
type CommentMap map[string]*Comment

// WithComment add a comment using the location and text information given in the CommentMap.
func WithComment(cm CommentMap) EncodeOption {
	return func(e *Encoder) error {
		commentMap := map[*Path]*Comment{}
		for k, v := range cm {
			path, err := PathString(k)
			if err != nil {
				return err
			}
			commentMap[path] = v
		}
		e.commentMap = commentMap
		return nil
	}
}

// CommentToMap apply the position and content of comments in a YAML document to a CommentMap.
func CommentToMap(cm CommentMap) DecodeOption {
	return func(d *Decoder) error {
		if cm == nil {
			return ErrInvalidCommentMapValue
		}
		d.toCommentMap = cm
		return nil
	}
}
