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

package generators

import (
	"errors"
	"fmt"
	"time"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/client"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/repository"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
)

type ServiceOperations interface {
	Generate(rqi interface{}) error
}

type ServiceOperationsImpl struct {
	ApiHubUrl     string
	ApiHubKey     string
	agentName     string
	kubeNameSpace string
	workSpace     string
	serviceId     string
	db            db.ConnectionProvider
	apihubClient  client.ApihubClient
}

// NewServiceOperationsReport
// creates a service operation report instance
func NewServiceOperationsReport(parameters ReportGeneratorParameters) (*ServiceOperationsImpl, error) {
	return &ServiceOperationsImpl{
		agentName:     parameters.AgentName,     //"k8s-apps3_api-hub-dev", // "k8sApps3-api-hub-dev",
		kubeNameSpace: parameters.KubeNameSpace, //"api-hub-dev",
		workSpace:     parameters.WorkSpace,     //"NC",
		serviceId:     view.EmptyString,
		db:            parameters.Db,
		apihubClient:  parameters.ApihubClient,
	}, nil
}

const (
	getStatusError = "unable to get report status for %s: %v"
	// ServiceVersionRequested - service version provided by user
	ServiceVersionRequested = "requested"
	// ServiceVersionRecent - the most service version received from APIHUB
	ServiceVersionRecent = "requested"
	// ServiceOperationPageSize limit service operations per page (used in APIHUB backend)
	ServiceOperationPageSize = 100
)

// Generate
// generates report according to the request
func (rep *ServiceOperationsImpl) Generate(rqi interface{}) error {
	rq := rqi.(view.ServiceReportRequest)
	contents, err := rep.apihubClient.GetPackagesVer(rep.apihubClient.GetSystemCtx(),
		view.PackagesSearchReq{
			ServiceName: rq.ServiceName,
			Kind:        "package",
		})
	if err != nil {
		return err
	}
	rt, err := repository.GetReportTypeByName(rep.db, entities.ServiceOperationReport)
	if err != nil {
		return err
	}
	rsc, err := repository.GetReportStatusByName(rep.db, entities.ReportStatusCreated)
	if err != nil {
		return fmt.Errorf(getStatusError, entities.ReportStatusCreated, err)
	}
	rsr, err := repository.GetReportStatusByName(rep.db, entities.ReportStatusReady)
	if err != nil {
		return fmt.Errorf(getStatusError, entities.ReportStatusReady, err)
	}
	rsf, err := repository.GetReportStatusByName(rep.db, entities.ReportStatusFailed)
	if err != nil {
		return fmt.Errorf(getStatusError, entities.ReportStatusFailed, err)
	}
	report := entities.ReportEntity{
		ReportId:       0,
		CreatedAt:      time.Now(),
		ReportTypeId:   rt.Id,
		ReportStatusId: rsc.Id,
		ReportUuid:     rq.ReportUuid,
	}
	if len(contents.Packages) < 1 {
		return fmt.Errorf("no packages found for service %s", rq.ServiceName)
	}
	rep.serviceId = contents.Packages[0].Id
	if rq.ServiceVersion == view.EmptyString {
		rq.ServiceVersion = contents.Packages[0].LastReleaseVersionDetails.Version
		rq.VersionStatus = ServiceVersionRecent // received from APIHUB
	} else {
		rq.VersionStatus = ServiceVersionRequested // user provided
	}
	err = repository.SetReportParameters(&report, rq)
	if err != nil {
		return err
	}
	err = repository.InsertReport(rep.db, &report)
	if err != nil {
		return err
	}
	err = rep.queryServiceOperations(rq, rep.serviceId, report.ReportId)
	if err == nil {
		report.ReportStatusId = rsr.Id
		log.Printf("report with id %d created", report.ReportId)
	} else {
		log.Warnf("report with id %d created with issue: %v", report.ReportId, err)
		report.ReportStatusId = rsf.Id
	}
	report.CompletedAt = time.Now()
	reportUpdateError := repository.UpdateReport(rep.db, report)
	if err == nil && reportUpdateError != nil {
		err = reportUpdateError
	}
	return err
}

// queryServiceOperations
// requests {API HUB URL}/api/v2/packages/{packageId}/versions/{version}/{apiType}/operations
// with parameters:
// serviceName A.K.A. packageId
// serviceVersion A.K.A. version
// "rest" A.K.A. apiType
func (rep *ServiceOperationsImpl) queryServiceOperations(rq view.ServiceReportRequest, serviceId string, reportId int) error {
	var err error
	// it is impossible to get all the operations at once - use paging
	currentPage := 0
	operationsOnPage := ServiceOperationPageSize
	opCount := 0
	cachedOpCount := 0
	// until the last (incomplete) page
	for operationsOnPage >= ServiceOperationPageSize {
		contents, errGetOps := rep.apihubClient.GetVersionRestOperationsWithData(
			rep.apihubClient.GetSystemCtx(), serviceId, rq.ServiceVersion, ServiceOperationPageSize, currentPage)
		if errGetOps != nil {
			return fmt.Errorf("unable to request service operations from APIHUB: %v", errGetOps)
		}
		if contents == nil {
			break
		}
		operationsOnPage = len(contents.Operations)
		currentPage++ // to the next page
		if operationsOnPage < 1 {
			break // no operations on page - break the loop
		}
		// dump service operation into DB, count occurrences, fill operation status
		tmpOptUpd := new(entities.ServicePacket)
		for _, op := range contents.Operations {
			opCount++
			tmpOpStat := entities.NewReportServiceOperation(reportId, op.OperationId, op.Path, op.Method, view.OperationNotFound)
			whereClause := "capture_id=? and request_method=? and "
			pathParam := tmpOpStat.Path
			if tmpOpStat.Regexp != view.EmptyString {
				whereClause += "not regexp_match(request_path, ?) is null "
				pathParam = tmpOpStat.Regexp
			} else {
				whereClause += "request_path = ? "
			}
			tmpOpStat.HitCount, err = rep.db.GetConnection().Model(tmpOptUpd).
				Where(whereClause, rq.CaptureId, tmpOpStat.Method, pathParam).Count()
			if err == nil {
				if tmpOpStat.HitCount > 0 {
					tmpOpStat.Status = view.OperationFound
				}
			} else {
				if !errors.Is(err, pg.ErrNoRows) {
					log.Debugf("unable to get hit count %s for service operation %s to report id %d: %v", whereClause, tmpOpStat.Path, reportId, err)
				}
			}
			err = entities.InsertReportServiceOperation(rep.db, &tmpOpStat)
			if err != nil {
				log.Debugf("unable to store service operation %s in Db: %v", op.Path, err)
			} else {
				cachedOpCount++
			}
		}
	}
	if opCount > 0 {
		if cachedOpCount != opCount {
			return fmt.Errorf("not all operations for report id %d were cached in Db: %v", reportId, err)
		}
		// collect affected packets
		_, err = rep.db.GetConnection().Exec(`
				insert into report_affected_rows (report_id, reference_id, reference_type, hit_count) 
				(select report_id, packet_id, ?, 1 from report_service_operations join service_packets 
				    on ((not regexp_match(request_path, report_service_operations.operation_path_re) is null) or request_path=report_service_operations.operation_path) and request_method=operation_method  
				where report_id=? and capture_id=?)`, //  ON CONFLICT (report_id, reference_id, reference_type) do nothing
			entities.ReportAffectedPacket, reportId, rq.CaptureId)
		if err != nil {
			return fmt.Errorf("unable to insert affected rows for packets in Db: %v", err)
		}
	}
	// insert packets which are not listed in service operations into output table
	sql3 := `insert into report_service_operations2
    	(report_id, src_peer, dst_peer, operation_title, operation_path, 
    	 operation_method, operation_status, hit_count)
	select ? as report_id, src_peer, dst_peer, '' as op_title, request_path, request_method, ? as op_status, sum(hit_count) as hit_count from 
		(select 
		case
			when length(coalesce(sas.service_name,''))<1 
			then concat(coalesce(sas.ip_address,''),':',to_char(source_port,'FM99999'))
			else sas.service_name end as src_peer,
		case
			when length(coalesce(sad.service_name, ''))<1 
			then concat(coalesce(sad.ip_address,''),':',to_char(source_port,'FM99999'))
			else sad.service_name end as dst_peer,
		request_path, 
		request_method, 
		1 as hit_count
		from service_packets rsp
		left join service_addresses sas on sas.address_id = source_id
		left join service_addresses sad on sad.address_id = dest_id
		where rsp.capture_id = ?
			and not exists (select null from report_affected_rows where report_id = ? and reference_id=packet_id and reference_type = ?)
			and not request_path is null) t2
		group by
			src_peer, dst_peer, request_path, request_method`
	_, err = rep.db.GetConnection().Exec(sql3, reportId, view.OperationExtra, rq.CaptureId, reportId, entities.ReportAffectedPacket)
	if err != nil {
		return fmt.Errorf("unable to insert operations not belong service in Db: %v", err)
	}
	// copy data
	// add previously collected operations
	sqlOp := `insert into report_service_operations2
    	(report_id, src_peer, dst_peer, operation_title, operation_path, 
    	 operation_method, operation_status, hit_count)
	select report_id, src_peer, dst_peer, operation_title, operation_path, operation_method, 
	       operation_status, sum(hit_count) as hit_count from (
		select report_id,	
		case
			when length(coalesce(sas.service_name,''))<1 
			then concat(coalesce(sas.ip_address,''),':',to_char(source_port,'FM99999'))
			else sas.service_name end as src_peer,
		case
			when length(coalesce(sad.service_name, ''))<1 
			then concat(coalesce(sad.ip_address,''),':',to_char(source_port,'FM99999'))
			else sad.service_name end as dst_peer,
		operation_title, operation_path, operation_method, operation_status,
		hit_count from 
	(select report_id, source_id, source_port, dest_id, dest_port,
		operation_title, operation_path, operation_method, operation_status,
		case when source_id is null then 0 else count(sp.*) end as hit_count
	from
		report_service_operations rps
	left JOIN service_packets sp
		on ((not regexp_match(request_path, rps.operation_path_re) is null)
			or request_path = rps.operation_path) 
			and rps.operation_method=sp.request_method
			and sp.capture_id = ?
	where
		rps.report_id = ?
	group by
	    report_id, source_id, source_port, dest_id, dest_port,
		operation_title, operation_path, operation_method, operation_status) rsp
	left join service_addresses sas on sas.address_id = rsp.source_id
	left join service_addresses sad on sad.address_id = rsp.dest_id) t2
	group by report_id, src_peer, dst_peer, operation_title, operation_path, 
	         operation_method, operation_status`
	_, err = rep.db.GetConnection().Exec(sqlOp, rq.CaptureId, reportId)
	if err != nil {
		if !errors.Is(err, pg.ErrNoRows) {
			log.Debugf("SQL:%s", sqlOp)
			log.Debugf("ERROR: %v", err)
			return fmt.Errorf("unable to select service operations into report: %v", err)
		}
	}
	// store report data
	reportData := make([]entities.ReportServiceOperationWithPeers, 0)
	// building query
	reportQuery2 := rep.db.GetConnection().Model(&reportData).Where("report_id=?", reportId).Order("operation_path", "operation_method", "src_peer", "dst_peer", "operation_title")
	if reportQuery2 == nil {
		return fmt.Errorf("no report data found for report id %d", reportId)
	}
	// execute query
	err = reportQuery2.Select()
	if err != nil {
		if !errors.Is(err, pg.ErrNoRows) {
			return err
		}
	}
	foundRows := 0
	insertedRows := 0
	// copy data
	foundRows, insertedRows, err = insertReportData(rep.db, reportData, reportId, rq.ServiceName)
	if insertedRows == foundRows {
		if log.GetLevel() != log.DebugLevel && log.GetLevel() != log.TraceLevel {
			// delete intermediate operation data if not being debugged
			intOperationData := new(entities.ReportServiceOperation)
			_, err = rep.db.GetConnection().Model(intOperationData).Where("report_id=?", reportId).Delete()
			if err != nil {
				return fmt.Errorf("unable to delete intermediate operation data for report id %d: %v", reportId, err)
			}
			// delete intermediate reference data
			intermediateRefData := new(entities.ReportAffectedRef)
			_, err = rep.db.GetConnection().Model(intermediateRefData).Where("report_id=?", reportId).Delete()
			if err != nil {
				return fmt.Errorf("unable to delete intermediate reference data for report id %d: %v", reportId, err)
			}
		}
	} else {
		return fmt.Errorf("%d rows inserted instead of %d for extras report id %d: %v", insertedRows, foundRows, reportId, err)
	}
	return nil
}

// insertReportData
// copy selected into report data table
func insertReportData(db db.ConnectionProvider, reportData []entities.ReportServiceOperationWithPeers, reportId int, serviceName string) (int, int, error) {
	var err error = nil
	insertedRows := 0
	foundRows := 0
	for _, reportOperationRow := range reportData {
		foundRows++
		senderService := reportOperationRow.Sender
		if len(senderService) < 2 {
			senderService = view.EmptyString
		}
		receiverService := reportOperationRow.Receiver
		if len(receiverService) < 2 {
			receiverService = serviceName
		}
		reportRowData := view.OperationStatusWithPeers{
			OperationStatus: view.OperationStatus{
				Id:       reportOperationRow.OperationId,
				Path:     reportOperationRow.Path,
				Method:   reportOperationRow.Method,
				Status:   reportOperationRow.Comment,
				HitCount: reportOperationRow.Occurrences,
				Peers:    nil,
			},
			Source:      senderService,
			Destination: receiverService,
		}
		reportRow := new(entities.ReportDataRow)
		reportRow.ReportId = reportId
		err = repository.SetReportRowData(reportRow, reportRowData)
		if err == nil {
			err = repository.InsertReportRow(db, reportRow)
			if err == nil {
				insertedRows++
			} else {
				log.Debugf("unable to insert report data for report id %d: %v", reportId, err)
			}
		} else {
			log.Debugf("unable to set report row data for report id %d: %v", reportId, err)
		}
	}
	return foundRows, insertedRows, err
}
