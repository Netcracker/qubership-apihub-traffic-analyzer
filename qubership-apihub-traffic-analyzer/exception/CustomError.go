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

package exception

import (
	"fmt"
	"strings"
)

type CustomError struct {
	Status  int                    `json:"status"`
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"message,omitempty"`
	Params  map[string]interface{} `json:"params,omitempty"`
	Debug   string                 `json:"debug,omitempty"`
}

func (c CustomError) Error() string {
	msg := c.Message
	for k, v := range c.Params {
		//todo make smart replace (e.g. now it replaces $projectId if we have $project in params)
		msg = strings.ReplaceAll(msg, "$"+k, fmt.Sprintf("%v", v))
	}
	if c.Debug != "" {
		return msg + " | " + c.Debug
	} else {
		return msg
	}
}

// todo replace with CustomError
type NotFoundError struct {
	Id      string
	Name    string
	Message string
}

func (g NotFoundError) Error() string {
	if g.Message != "" {
		return g.Message
	}
	if g.Id != "" {
		return fmt.Sprintf("entity with id = %s not found", g.Id)
	} else {
		return fmt.Sprintf("entity with name = %s not found", g.Name)
	}
}

// todo replace with CustomError
type ContentNotFoundError struct {
	ContentId string
	Slug      string
}

func (c ContentNotFoundError) Error() string {
	if c.ContentId != "" {
		return fmt.Sprintf("content with contentId = %v not found", c.ContentId)
	} else {
		return fmt.Sprintf("content with slug = %v not found", c.Slug)
	}
}
