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

package secctx

import (
	"github.com/shaj13/go-guardian/v2/auth"
	"net/http"
	"strings"
)

type SecurityContext interface {
	GetUserId() string
	GetUserToken() string
	GetApiKey() string
	IsSystem() bool
}

func Create(r *http.Request) SecurityContext {
	token := getAuthorizationToken(r)
	if token != "" {
		user := auth.User(r)
		userId := user.GetID()
		return &securityContextImpl{
			userId:   userId,
			token:    token,
			apiKey:   "",
			isSystem: false,
		}
	} else {
		apiKey := getApihubApiKey(r)
		userId := "api-key_" + apiKey
		return &securityContextImpl{
			userId:   userId,
			token:    "",
			apiKey:   apiKey,
			isSystem: false,
		}
	}
}

func CreateSystemContext() SecurityContext {
	return &securityContextImpl{isSystem: true}
}

type securityContextImpl struct {
	userId   string
	token    string
	apiKey   string
	isSystem bool
}

func getAuthorizationToken(r *http.Request) string {
	authorizationHeaderValue := r.Header.Get("authorization")
	return strings.ReplaceAll(authorizationHeaderValue, "Bearer ", "")
}

func getApihubApiKey(r *http.Request) string {
	return r.Header.Get("api-key")
}

func (ctx securityContextImpl) GetUserId() string {
	return ctx.userId
}
func (ctx securityContextImpl) GetUserToken() string {
	return ctx.token
}
func (ctx securityContextImpl) GetApiKey() string {
	return ctx.apiKey
}
func (ctx securityContextImpl) IsSystem() bool { return ctx.isSystem }
