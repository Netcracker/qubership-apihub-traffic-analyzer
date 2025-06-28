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
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/openapi/yaml/ast"
	"golang.org/x/xerrors"
)

var (
	ErrInvalidQuery               = xerrors.New("invalid query")
	ErrInvalidPath                = xerrors.New("invalid path instance")
	ErrInvalidPathString          = xerrors.New("invalid path string")
	ErrNotFoundNode               = xerrors.New("node not found")
	ErrUnknownCommentPositionType = xerrors.New("unknown comment position type")
	ErrInvalidCommentMapValue     = xerrors.New("invalid comment map value. it must be not nil value")
)

func ErrUnsupportedHeadPositionType(node ast.Node) error {
	return xerrors.Errorf("unsupported comment head position for %s", node.Type())
}

// IsInvalidQueryError whether err is ErrInvalidQuery or not.
func IsInvalidQueryError(err error) bool {
	return xerrors.Is(err, ErrInvalidQuery)
}

// IsInvalidPathError whether err is ErrInvalidPath or not.
func IsInvalidPathError(err error) bool {
	return xerrors.Is(err, ErrInvalidPath)
}

// IsInvalidPathStringError whether err is ErrInvalidPathString or not.
func IsInvalidPathStringError(err error) bool {
	return xerrors.Is(err, ErrInvalidPathString)
}

// IsNotFoundNodeError whether err is ErrNotFoundNode or not.
func IsNotFoundNodeError(err error) bool {
	return xerrors.Is(err, ErrNotFoundNode)
}

// IsInvalidTokenTypeError whether err is ast.ErrInvalidTokenType or not.
func IsInvalidTokenTypeError(err error) bool {
	return xerrors.Is(err, ast.ErrInvalidTokenType)
}

// IsInvalidAnchorNameError whether err is ast.ErrInvalidAnchorName or not.
func IsInvalidAnchorNameError(err error) bool {
	return xerrors.Is(err, ast.ErrInvalidAnchorName)
}

// IsInvalidAliasNameError whether err is ast.ErrInvalidAliasName or not.
func IsInvalidAliasNameError(err error) bool {
	return xerrors.Is(err, ast.ErrInvalidAliasName)
}
