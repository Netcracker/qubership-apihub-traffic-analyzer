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

package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/exception"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/secctx"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

const (
	packagesByService = "%s/api/v2/packages"
	operationsUri     = "%s/api/v2/packages/%s/versions/%s/%s/operations"
	restOperationType = "rest"
	operationKind     = "kind"
	filterLimit       = "limit"
	filterPage        = "page"
)

type apihubClientImpl struct {
	apihubUrl   string
	accessToken string
	apiHubHost  string
	SystemCtx   secctx.SecurityContext
}

// ApihubClient
// public interface
type ApihubClient interface {
	GetVersionRestOperationsWithData(ctx secctx.SecurityContext, packageId, version string, limit, page int) (*view.RestOperations, error)
	GetPackagesVer(ctx secctx.SecurityContext, searchReq view.PackagesSearchReq) (*view.Packages, error)
	GetPackages(ctx secctx.SecurityContext, searchReq view.PackagesSearchReq) (*view.SimplePackages, error)
	GetSystemCtx() secctx.SecurityContext
}

// NewApihubClient
// creates an instance of the client
func NewApihubClient(apihubUrl, accessToken string) ApihubClient {
	parsedApihubUrl, err := url.Parse(apihubUrl)
	apihubHost := view.EmptyString
	if err != nil {
		log.Errorf("Can't parse apihub url: %v", err)
	} else {
		apihubHost = parsedApihubUrl.Hostname()
	}
	log.Printf("apihub client created for URL  %s", apihubUrl)
	log.Debugf("apihub client uses apihubHost  %s", apihubHost)
	log.Debugf("apihub client uses accessToken %s", accessToken)
	return &apihubClientImpl{apihubUrl: apihubUrl, accessToken: accessToken, apiHubHost: apihubHost, SystemCtx: secctx.CreateSystemContext()}
}

// GetVersionRestOperationsWithData
// get REST operations for package and version
func (a apihubClientImpl) GetVersionRestOperationsWithData(ctx secctx.SecurityContext, packageId, version string, limit, page int) (*view.RestOperations, error) {
	req := makeRequest(ctx, a.accessToken, a.apiHubHost)
	req.SetQueryParam("includeData", "true")
	req.SetQueryParam(filterLimit, fmt.Sprint(limit))
	req.SetQueryParam(filterPage, fmt.Sprint(page))
	resp, err := req.Get(fmt.Sprintf(operationsUri,
		a.apihubUrl,
		url.PathEscape(packageId),
		url.PathEscape(version),
		restOperationType))
	if err != nil {
		return nil, fmt.Errorf("failed to get version rest operations. Error - %s", err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusNotFound {
			return nil, nil
		}
		if authErr := checkUnauthorized(resp); authErr != nil {
			return nil, authErr
		}
		return nil, fmt.Errorf("failed to get version rest operations: status code %d %v", resp.StatusCode(), err)
	}

	var restOperations view.RestOperations
	err = json.Unmarshal(resp.Body(), &restOperations)
	if err != nil {
		return nil, err
	}
	return &restOperations, nil
}

// GetPackagesVer
// returns package data with version
func (a apihubClientImpl) GetPackagesVer(ctx secctx.SecurityContext, searchReq view.PackagesSearchReq) (*view.Packages, error) {
	req := makeRequest(ctx, a.accessToken, a.apiHubHost)

	uri := fmt.Sprintf(packagesByService, a.apihubUrl)
	if searchReq.Limit == 0 {
		searchReq.Limit = 100
	}
	if searchReq.TextFilter != view.EmptyString {
		req.SetQueryParam("textFilter", searchReq.TextFilter)
	}
	if searchReq.ServiceName != view.EmptyString {
		req.SetQueryParam("serviceName", searchReq.ServiceName)
	}
	if searchReq.ParentId != view.EmptyString {
		req.SetQueryParam("parentId", searchReq.ParentId)
	}
	if searchReq.ShowAllDescendants {
		req.SetQueryParam("showAllDescendants", strconv.FormatBool(searchReq.ShowAllDescendants))
	}
	if searchReq.Kind != view.EmptyString {
		req.SetQueryParam(operationKind, searchReq.Kind)
	}
	if searchReq.Limit > 0 && searchReq.Page > 0 {
		req.SetQueryParam(filterLimit, fmt.Sprint(searchReq.Limit))
		req.SetQueryParam(filterPage, fmt.Sprint(searchReq.Page))
	}
	req.SetQueryParam("lastReleaseVersionDetails", "true")
	resp, err := req.Get(uri)
	statusCode := http.StatusInternalServerError
	if resp != nil {
		statusCode = resp.StatusCode()
	}
	if err != nil {
		return nil,
			fmt.Errorf("failed to get package by serarchReq : %v. Response status code %d %v",
				searchReq, statusCode, err)
	}
	if statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			return nil, nil
		}
		if authErr := checkUnauthorized(resp); authErr != nil {
			return nil, authErr
		}
		return nil,
			fmt.Errorf("failed to get packages by request -  %v : status code %d %v",
				searchReq, statusCode, err)
	}

	var packages view.Packages

	err = json.Unmarshal(resp.Body(), &packages)
	if err != nil {
		return nil, err
	}
	if len(packages.Packages) == 0 {
		return nil, nil
	}
	return &packages, nil
}

// GetPackages
// get package list by filter
func (a apihubClientImpl) GetPackages(ctx secctx.SecurityContext, searchReq view.PackagesSearchReq) (*view.SimplePackages, error) {
	req := makeRequest(ctx, a.accessToken, a.apiHubHost)

	uri := fmt.Sprintf(packagesByService, a.apihubUrl)
	if searchReq.Limit == 0 {
		searchReq.Limit = 100
	}
	if searchReq.TextFilter != "" {
		req.SetQueryParam("textFilter", searchReq.TextFilter)
	}
	if searchReq.ServiceName != "" {
		req.SetQueryParam("serviceName", searchReq.ServiceName)
	}
	if searchReq.ParentId != "" {
		req.SetQueryParam("parentId", searchReq.ParentId)
	}
	if searchReq.ShowAllDescendants {
		req.SetQueryParam("showAllDescendants", strconv.FormatBool(searchReq.ShowAllDescendants))
	}
	if searchReq.Kind != "" {
		req.SetQueryParam(operationKind, searchReq.Kind)
	}
	if searchReq.Limit > 0 && searchReq.Page > 0 {
		req.SetQueryParam(filterLimit, fmt.Sprint(searchReq.Limit))
		req.SetQueryParam(filterPage, fmt.Sprint(searchReq.Page))
	}
	resp, err := req.Get(uri)
	statusCode := http.StatusInternalServerError
	if resp != nil {
		statusCode = resp.StatusCode()
	}
	if err != nil {
		return nil,
			fmt.Errorf("failed to get package by serarchReq : %v. Response status code %d %v",
				searchReq, statusCode, err)
	}
	if statusCode != http.StatusOK {
		if statusCode == http.StatusNotFound {
			return nil, nil
		}
		if authErr := checkUnauthorized(resp); authErr != nil {
			return nil, authErr
		}
		return nil,
			fmt.Errorf("failed to get packages by request -  %v : status code %d %v",
				searchReq, statusCode, err)
	}

	var packages view.SimplePackages

	err = json.Unmarshal(resp.Body(), &packages)
	if err != nil {
		return nil, err
	}
	if len(packages.Packages) == 0 {
		return nil, nil
	}
	return &packages, nil
}

// GetSystemCtx
// return internal (default) security context
func (a apihubClientImpl) GetSystemCtx() secctx.SecurityContext {
	return a.SystemCtx
}

// checkUnauthorized
// validate authorization
func checkUnauthorized(resp *resty.Response) error {
	if resp != nil &&
		(resp.StatusCode() == http.StatusUnauthorized || resp.StatusCode() == http.StatusForbidden) {
		log.Errorf("Incorrect api key detected!")
		return &exception.CustomError{
			Status:  http.StatusFailedDependency,
			Code:    exception.NoApihubAccess,
			Message: exception.NoApihubAccessMsg,
			Params:  map[string]interface{}{"code": strconv.Itoa(resp.StatusCode())},
		}
	}
	return nil
}

// makeRequest
// makes an API hub request
func makeRequest(ctx secctx.SecurityContext, apiHubKey, serverName string) *resty.Request {
	tr := http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	cl := http.Client{Transport: &tr, Timeout: time.Second * 60}

	client := resty.NewWithClient(&cl)
	client.SetRedirectPolicy(resty.DomainCheckRedirectPolicy(serverName))
	req := client.R()
	if ctx.IsSystem() {
		if apiHubKey != view.EmptyString {
			req.SetHeader(view.ApiKeyHeader, apiHubKey)
		}
	} else {
		if ctx.GetUserToken() != "" {
			req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", ctx.GetUserToken()))
		} else if ctx.GetApiKey() != view.EmptyString {
			if apiHubKey != view.EmptyString {
				req.SetHeader(view.ApiKeyHeader, apiHubKey)
			}
		}

	}
	return req
}

// checkAndGetBody
// reads body from "raw" HTTP response
//func checkAndGetBody(r *http.Response) ([]byte, error) {
//	if r == nil {
//		return nil, fmt.Errorf("nil response")
//	}
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Debugf("unable to defer response body: %v", err)
//		}
//	}(r.Body)
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		return nil, err
//	}
//	return body, nil
//}
