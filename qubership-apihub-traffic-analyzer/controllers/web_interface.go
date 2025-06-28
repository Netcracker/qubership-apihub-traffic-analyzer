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

package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/client"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/exception"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/readers"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/reports/generators"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/reports/renderers"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/repository"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/service"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/utils"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	OnCaptureCleanup(w http.ResponseWriter, r *http.Request)
	OnCaptureDelete(w http.ResponseWriter, r *http.Request)
	OnCaptureLoad(w http.ResponseWriter, r *http.Request)
	OnCaptureLoadStatus(w http.ResponseWriter, r *http.Request)
	OnStatus(w http.ResponseWriter, r *http.Request)
	Shutdown()
	OnServiceOperationsReportGenerate(w http.ResponseWriter, r *http.Request)
	OnServiceOperationsReportOutput(w http.ResponseWriter, r *http.Request)
}

const (
	HttpContentType        = "Content-Type"
	HttpContentJson        = "application/json"
	HttpContentOctetStream = "application/octet-stream"
	invalidApiKey          = "API key not match"
	emptyApiKey            = "empty API key not allowed in production mode"
	emptyCaptureId         = "Capture Id is empty"
	requestBodyDeferError  = "unable to defer request body. error: %v"
	StopAsync              = "STOP"
)

type webService struct {
	entities.WebServiceConfig
	Headers       repository.HttpHeadersCache
	Packets       repository.PacketCache
	Peers         repository.ServiceAddressRepository
	history       map[string]view.LoadHistoryValue
	reports       chan string
	s3            service.CloudStorage
	db            db.ConnectionProvider
	mapLock       sync.Mutex
	kubeNameSpace string
	workSpace     string
	apihubClient  client.ApihubClient
}

// NewService
// creates a new web service instance
func NewService(cfg entities.WebServiceConfig,
	headers repository.HttpHeadersCache,
	packets repository.PacketCache,
	peers repository.ServiceAddressRepository,
	s3 service.CloudStorage,
	pdb db.ConnectionProvider,
	kubeNameSpace,
	workSpace string,
	apihubClient client.ApihubClient) Service {
	ws := &webService{
		WebServiceConfig: cfg,
		Headers:          headers,
		Packets:          packets,
		Peers:            peers,
		history:          make(map[string]view.LoadHistoryValue),
		reports:          make(chan string),
		s3:               s3,
		db:               pdb,
		mapLock:          sync.Mutex{},
		kubeNameSpace:    kubeNameSpace,
		workSpace:        workSpace,
		apihubClient:     apihubClient,
	}
	utils.SafeAsync(func() {
		for {
			captureId := <-ws.reports
			if captureId == StopAsync {
				break
			}
			val, found := ws.history[captureId]
			if found {
				if completed, err := view.HistoryRecordCompleted(val); completed {
					if err == nil {
						log.Printf("capture '%s' loaded successfully", captureId)
					} else {
						log.Errorf("capture '%s' load finished with error: %v", captureId, err)
					}
				}
				ws.mapLock.Lock()
				val.EndDateTime = view.GetHistoryDateTimeString()
				ws.history[captureId] = val
				ws.mapLock.Unlock()
			} else {
				log.Errorf("unexpected capture id %s", captureId)
			}
		}
	})
	return ws
}

// RespondWithJson
// respond to the other side with custom JSON
func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set(HttpContentType, HttpContentJson)
	w.WriteHeader(code)
	write, err := w.Write(response)
	if err != nil {
		log.Debugf("%d response bytes written with error: %v", write, err)
	}
}

// RespondWithCustomError
// respond to the other side with custom (error) text
func RespondWithCustomError(w http.ResponseWriter, err *exception.CustomError) {
	log.Debugf("Request failed. Code = %d. Message = %s. Params: %v. Debug: %s", err.Status, err.Message, err.Params, err.Debug)
	RespondWithJson(w, err.Status, err)
}

// getStringParam
// get string value from path by its name
func getStringParam(r *http.Request, paramName string) string {
	params := mux.Vars(r)
	return params[paramName]
}

// setCaptureError
// sets capture status and error, creates record if it does not exist
func (ws *webService) setCaptureError(captureId string, err error) {
	ws.mapLock.Lock()
	defer ws.mapLock.Unlock()
	val, found := ws.history[captureId]
	if found {
		val.Error = err
		ws.history[captureId] = val
	} else {
		ws.history[captureId] = view.LoadHistoryValue{
			BeginDateTime: view.GetHistoryDateTimeString(),
			EndDateTime:   view.EmptyString,
			Error:         err,
		}
	}
}

// OnCaptureLoad
// serves the capture load requests
func (ws *webService) OnCaptureLoad(w http.ResponseWriter, r *http.Request) {
	_, err := ws.checkAndGetBody(w, r)
	if err != nil {
		return
	}
	captureId := getStringParam(r, view.CaptureIdParam)
	if captureId == view.EmptyString {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.ContentIdNotFound,
			Message: exception.ContentIdNotFoundMsg,
			Debug:   emptyCaptureId,
		})
		return
	}
	if captureId == StopAsync {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.InvalidParameterValue,
			Message: exception.InvalidParameterValueMsg,
			Debug:   StopAsync,
		})
		return
	}
	status, found := ws.history[captureId]
	if found {
		completed, _ := view.HistoryRecordCompleted(status)
		if completed {
			// if completed - delete it for further re-creation
			ws.mapLock.Lock()
			delete(ws.history, captureId)
			ws.mapLock.Unlock()
		} else {
			// incomplete capture
			if status.Error == nil {
				// indicate an incomplete state
				RespondWithJson(w, http.StatusPartialContent, view.EmptyString)
			} else {
				// indicate an incomplete state with error
				RespondWithJson(w, http.StatusPartialContent, status.Error.Error())
			}
			return // avoid to open concurrent captures
		}
	}
	ws.setCaptureError(captureId, nil)
	utils.SafeAsync(func() {
		log.Printf("starting process files for capture %s", captureId)
		rdr := readers.NewCaptureReader(ws.Headers, ws.Packets, ws.Peers, ws.db, ws.WorkDir)
		_, err := ws.s3.ProcessCaptureFiles(captureId, rdr)
		ws.setCaptureError(captureId, err)
		err = rdr.Close()
		if err != nil {
			log.Warnf("unable to close reader: %v", err)
		}
		ws.reports <- captureId
	})
	RespondWithJson(w, http.StatusAccepted, "loading")
}

// OnCaptureLoadStatus
// returns status for a particular capture ID
func (ws *webService) OnCaptureLoadStatus(w http.ResponseWriter, r *http.Request) {
	_, err := ws.checkAndGetBody(w, r)
	if err != nil {
		return
	}
	captureId := getStringParam(r, view.CaptureIdParam)
	if captureId == "" {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.ContentIdNotFound,
			Message: exception.ContentIdNotFoundMsg,
			Debug:   emptyCaptureId,
		})
		return
	}
	if captureId == StopAsync {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.InvalidParameterValue,
			Message: exception.InvalidParameterValueMsg,
			Debug:   StopAsync,
		})
		return
	}
	status, found := ws.history[captureId]
	if !found {
		RespondWithJson(w, http.StatusNotFound, fmt.Sprintf("capture '%s' was not found", captureId))
	} else {
		completed, err := view.HistoryRecordCompleted(status)
		if completed {
			if err == nil {
				RespondWithJson(w, http.StatusOK, fmt.Sprintf("capture '%s' was loaded at %s", captureId, status.EndDateTime))
			} else {
				RespondWithJson(w, http.StatusExpectationFailed, fmt.Sprintf("capture '%s' was failed at %s : %v", captureId, status.EndDateTime, err))
			}
		} else {
			RespondWithJson(w, http.StatusCreated, fmt.Sprintf("capture '%s' is still loading", captureId))
		}
	}
}

// OnStatus
// responds to a cloud status requests (/live, /ready, /startup)
func (ws *webService) OnStatus(w http.ResponseWriter, _ *http.Request) {
	RespondWithJson(w, http.StatusOK, view.GetHistoryDateTimeString()) // always respond OK to calm the watchdogs
}

// checkAndGetBody
// checks API key and reads body contents
func (ws *webService) checkAndGetBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if ws.APIkey != view.EmptyString {
		apiKeyHeader := r.Header.Get(view.ApiKeyHeader)
		if apiKeyHeader != ws.APIkey {
			RespondWithCustomError(w, &exception.CustomError{
				Status:  http.StatusUnauthorized,
				Code:    exception.ApiKeyNotFound,
				Message: exception.ApiKeyNotFoundMsg,
				Debug:   invalidApiKey,
			})
			return nil, fmt.Errorf(invalidApiKey)
		}
	} else {
		if ws.ProductionMode {
			RespondWithCustomError(w, &exception.CustomError{
				Status:  http.StatusUnauthorized,
				Code:    exception.EmptyParameter,
				Message: exception.EmptyParameterMsg,
				Debug:   emptyApiKey,
			})
			return nil, fmt.Errorf(emptyApiKey)
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Debugf(requestBodyDeferError, err)
		}
	}(r.Body)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.BadRequestBody,
			Message: exception.BadRequestBodyMsg,
			Debug:   err.Error(),
		})
		return nil, err
	}
	return body, nil
}

// Shutdown
// tries to perform a graceful service shutdown
func (ws *webService) Shutdown() {
	ws.reports <- StopAsync
}

// OnServiceOperationsReportGenerate
// generates report data for service operation report with capture id and service name/version
func (ws *webService) OnServiceOperationsReportGenerate(w http.ResponseWriter, r *http.Request) {
	body, err := ws.checkAndGetBody(w, r)
	if err != nil {
		return
	}
	var req view.ServiceReportRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.BadRequestBody,
			Message: exception.BadRequestBodyMsg,
			Debug:   err.Error(),
		})
		return
	}
	err = view.ValidateServiceReportRequest(req)
	if err != nil {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.EmptyParameter,
			Message: exception.EmptyParameterMsg,
			Debug:   err.Error(),
		})
		return
	}
	log.Printf("%v", req)
	uuid := utils.MakeUniqueId()
	req.ReportUuid = uuid
	rep, err := generators.NewReportGenerator(generators.ReportGeneratorParameters{
		ApihubClient:  ws.apihubClient,
		KubeNameSpace: ws.kubeNameSpace,
		WorkSpace:     ws.workSpace,
		WorkDir:       ws.WorkDir,
		AgentName:     ws.AgentName,
		ReportType:    generators.ServiceOperationReport,
		Db:            ws.db,
	})
	if err != nil {
		log.Errorf("error instantiating service operations report: %v", err)
	} else {
		utils.SafeAsync(func() {
			if req.ReportUuid != uuid {
				log.Debugf("service operations report: %s %s", req.ReportUuid, uuid)
				req.ReportUuid = uuid
			}
			err = rep.Generate(req)
			if err != nil {
				log.Warnf("unable to generate service operations report: %v", err)
			}
		})
	}
	RespondWithJson(w, http.StatusAccepted, view.ReportDataRequest{Id: uuid})
}

// OnServiceOperationsReportOutput
// makes a report data render and send it back if not timed out
func (ws *webService) OnServiceOperationsReportOutput(w http.ResponseWriter, r *http.Request) {
	body, err := ws.checkAndGetBody(w, r)
	if err != nil {
		return
	}
	var req view.ReportDataRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.BadRequestBody,
			Message: exception.BadRequestBodyMsg,
			Debug:   err.Error(),
		})
		return
	}
	err = view.ValidateReportDataRequest(&req)
	if err != nil {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.EmptyParameter,
			Message: exception.EmptyParameterMsg,
			Debug:   err.Error(),
		})
		return
	}
	var repRender renderers.ReportRenderer
	repRender, err = renderers.NewReportRenderer(ws.db, req, ws.WorkDir, generators.ServiceOperationReport)
	if err == nil {
		asyncChan := make(chan string)
		// render report asynchronously
		utils.SafeAsync(func() {
			// read report data to send it out
			// make contents
			err = repRender.MakeReportHeader()
			if err != nil {
				log.Debugf("unable to make report header: %v", err)
			}
			err = repRender.ProcessRows()
			if err != nil {
				log.Debugf("unable to process report rows: %v", err)
			}
			err = repRender.MakeReportFooter()
			if err != nil {
				log.Debugf("unable to make report footer: %v", err)
			}
			asyncChan <- view.EmptyString
		})
		select {
		case _ = <-asyncChan:
			break
		case <-time.After(120 * time.Second):
			{
				RespondWithCustomError(w, &exception.CustomError{
					Status:  http.StatusRequestTimeout,
					Code:    exception.ReportGenerationTimeOut,
					Message: exception.ReportGenerationTimeOutMsg,
					Debug:   "",
				})
				return
			}
		}
		switch req.Format {
		case view.ReportFormatExcel:
			{
				w.Header().Set(HttpContentType, HttpContentOctetStream)
				w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, repRender.GetFileName()))
				w.Header().Set("Content-Transfer-Encoding", "binary")
				w.Header().Set("Expires", "0")
			}
		case view.ReportFormatJson:
			w.Header().Set(HttpContentType, HttpContentJson)
		}
		w.WriteHeader(http.StatusOK)
		err = repRender.FlushData(w)
		if err != nil {
			log.Debugf("unable to pass report contents: %v", err)
		}
		repRender.Dispose()
	} else {
		if err != nil {
			RespondWithCustomError(w, &exception.CustomError{
				Status:  http.StatusNotFound,
				Code:    exception.InvalidParameterValue,
				Message: exception.InvalidParameterValueMsg,
				Debug:   err.Error(),
			})
			return
		}
	}
}

// OnCaptureDelete
// tries to delete capture data from S3/Minio
func (ws *webService) OnCaptureDelete(w http.ResponseWriter, r *http.Request) {
	_, err := ws.checkAndGetBody(w, r)
	if err != nil {
		return
	}
	captureId := getStringParam(r, view.CaptureIdParam)
	if captureId == view.EmptyString {
		RespondWithCustomError(w, &exception.CustomError{
			Status:  http.StatusBadRequest,
			Code:    exception.ContentIdNotFound,
			Message: exception.ContentIdNotFoundMsg,
			Debug:   emptyCaptureId,
		})
		return
	}
	status, found := ws.history[captureId]
	if found {
		completed, _ := view.HistoryRecordCompleted(status)
		if !completed {
			// incomplete (in progress) capture
			if status.Error == nil {
				// indicates an incomplete state
				RespondWithJson(w, http.StatusPartialContent, view.EmptyString)
			}
			return // avoid to open concurrent captures
		}
	}
	utils.SafeAsync(func() {
		log.Printf("trying to delete files for capture %s", captureId)
		deletedCount, err := ws.s3.DeleteCaptureFiles(captureId)
		if err == nil {
			if deletedCount > 0 {
				log.Printf("successfully deleted %d files for capture %s", deletedCount, captureId)
			} else {
				log.Printf("no files found for capture %s", captureId)
			}
		} else {
			log.Errorf("unable to delete files for capture %s: %v", captureId, err)
		}
	})
	RespondWithJson(w, http.StatusAccepted, "deleting")
}

// OnCaptureCleanup
// tries to delete capture data from S3/Minio
func (ws *webService) OnCaptureCleanup(w http.ResponseWriter, r *http.Request) {
	_, err := ws.checkAndGetBody(w, r)
	if err != nil {
		return
	}
	utils.SafeAsync(func() {
		log.Printf("trying to clean the capture files")
		deletedCount, err := ws.s3.CleanupCaptureFiles()
		if err == nil {
			if deletedCount > 0 {
				log.Printf("successfully cleaned %d files", deletedCount)
			} else {
				log.Printf("no files found to cleanup")
			}
		} else {
			log.Errorf("unable to cleanup capture files: %v", err)
		}
	})
	RespondWithJson(w, http.StatusAccepted, "cleaning up")
}
