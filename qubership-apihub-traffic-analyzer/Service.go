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

package main

import (
	"flag"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/client"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/controllers"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/db"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/entities"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/readers"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/reports/generators"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/repository"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/service"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/utils"
	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func makeServer(systemInfoService service.SystemInfoService, r *mux.Router) *http.Server {
	listenAddr := systemInfoService.GetListenAddress()

	log.Infof("Listen addr = %s", listenAddr)

	var corsOptions []handlers.CORSOption

	corsOptions = append(corsOptions,
		handlers.AllowedHeaders([]string{
			"Connection",
			"Accept-Encoding",
			"Content-Encoding",
			"X-Requested-With",
			controllers.HttpContentType,
			"Authorization"}))

	allowedOrigin := systemInfoService.GetOriginAllowed()
	if allowedOrigin != "" {
		corsOptions = append(corsOptions, handlers.AllowedOrigins([]string{allowedOrigin}))
	}
	corsOptions = append(corsOptions, handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodDelete}))

	return &http.Server{
		Handler:      handlers.CompressHandler(handlers.CORS(corsOptions...)(r)),
		Addr:         listenAddr,
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
}

func init() {
	// log rotation + log path
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = "."
	}
	mw := io.MultiWriter(os.Stderr, &lumberjack.Logger{
		Filename: basePath + "/logs/traffic-analyzer.log",
		MaxSize:  10, // megabytes
	})
	// log formatter
	//log.SetFormatter(&prefixed.TextFormatter{
	//	DisableColors:   true,
	//	TimestampFormat: "2006-01-02 15:04:05",
	//	FullTimestamp:   true,
	//	ForceFormatting: true,
	//})
	// log level
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel // "INFO"
	}
	log.SetLevel(logLevel)
	// apply log rotation
	log.SetOutput(mw)
}

func main() {
	var (
		reportName     string
		captureId      string
		serviceName    string
		serviceVersion string
		logLevel       string
	)
	sysInfo, err := service.NewSystemInfoService()
	err = sysInfo.Init()
	if err != nil {
		log.Fatalf("unable to initialize from configuration. Error: %v", err)
	} else {
		var (
			workDir   string
			baseDir   string
			connAttrs view.DbCredentials
		)
		flag.StringVar(&captureId, "capture-id", sysInfo.GetCaptureId(), "Capturing ID to aggregate")
		flag.StringVar(&workDir, "work-dir", sysInfo.GetWorkDir(), "Working directory for intermediate files")
		flag.StringVar(&baseDir, "base-dir", sysInfo.GetBasePath(), "Base directory for migration")
		flag.StringVar(&connAttrs.Host, "host", sysInfo.GetPGHost(), "DB server host")
		flag.StringVar(&connAttrs.Database, "instance", sysInfo.GetPGDB(), "DB instance name")
		flag.StringVar(&connAttrs.Password, "password", sysInfo.GetPGPassword(), "DB user password")
		flag.StringVar(&connAttrs.Username, "user", sysInfo.GetPGUser(), "DB user name")
		flag.StringVar(&connAttrs.Schema, "schema", view.EmptyString, "DB schema name")
		flag.StringVar(&connAttrs.SSLMode, "ssl-mode", sysInfo.GetPGSSLMode(), "SSL mode")
		flag.IntVar(&connAttrs.Port, "port", sysInfo.GetPGPort(), "DB server port")
		flag.StringVar(&reportName, "report-name", view.EmptyString, "report name to generate")
		flag.StringVar(&serviceName, "service-name", view.EmptyString, "service name to generate report")
		flag.StringVar(&serviceVersion, "service-version", view.EmptyString, "service version to generate report")
		flag.StringVar(&logLevel, "log-level", "info", "A logging level: (trace, debug, info, warning, error, fatal, panic)")
		flag.Parse()
		sysInfo.CmdLineOverride(captureId, baseDir, workDir, connAttrs)
	}
	err = sysInfo.Validated()
	if err != nil {
		log.Fatalf("configuration not valid: %v", err)
		return
	}
	// connection provider
	pdb := db.NewConnectionProvider(sysInfo.GetCredsFromEnv())
	// migration
	migrationResult := make(chan int, 1)
	dbMigrationService, err := service.NewDBMigrationService(pdb, sysInfo)
	if err != nil {
		log.Fatalf("Failed create dbMigrationService: " + err.Error())
	}
	go func() { // Do not use safe async here to enable panic
		_, _, _, err := dbMigrationService.Migrate(sysInfo.GetBasePath())
		if err != nil {
			log.Fatalf("Failed perform DB migration: " + err.Error())
			migrationResult <- -2
			time.Sleep(time.Second * 10) // Give a chance to read the unrecoverable error
		} else {
			migrationResult <- 1
		}
	}()
	// wait for migration to end
	reqCompleted := 0
	select {
	case res := <-migrationResult:
		reqCompleted = res // confirmed
	case <-time.After(time.Second * 5):
		reqCompleted = -1 // timeout reading channel
	}
	if reqCompleted != 1 {
		switch reqCompleted {
		case 0:
			log.Fatalf("unexpected migration result")
		case -1:
			log.Fatalf("timeout waiting for migration result")
		case -2:
			log.Fatalf("migration caused error")
		}
		return
	}
	capId := sysInfo.GetCaptureId()
	headersCache := repository.NewHttpHeadersCache(pdb)
	peersCache := repository.NewPeersCache(pdb)
	packetCache := repository.NewPacketCache(pdb, peersCache, headersCache)
	minioCfg := sysInfo.GetMinioStorageCreds()
	var s3 service.CloudStorage = nil
	s3, err = service.NewCloudStorage(*minioCfg)
	if err != nil {
		log.Warnf("unable to initialise cloud storage interface: %v", err)
	}
	// API hub client
	apihubClient := client.NewApihubClient(sysInfo.GetApiHubUrl(), sysInfo.GetApiHubAccessToken())
	switch strings.ToUpper(logLevel) {
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARNING":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
	if reportName != view.EmptyString {
		switch reportName {
		case view.ServiceOperationsReportPath:
			{
				rep, err := generators.NewServiceOperationsReport(generators.ReportGeneratorParameters{
					ApihubClient:  apihubClient,
					KubeNameSpace: sysInfo.GetNamespace(),
					WorkSpace:     sysInfo.GetWorkspace(),
					WorkDir:       sysInfo.GetWorkDir(),
					AgentName:     sysInfo.GetAgentName(),
					ReportType:    generators.ServiceOperationReport,
					Db:            pdb,
				})
				if err != nil {
					log.Fatalf("error creating service operations report - %v", err)
				}
				rq := view.ServiceReportRequest{
					CaptureId:      captureId,
					ServiceName:    serviceName,
					ServiceVersion: serviceVersion,
					ReportUuid:     utils.MakeUniqueId(),
				}
				err = rep.Generate(rq)
				if err != nil {
					log.Fatalf("unable to generate service operations report - %v", err)
				}
			}
		default:
			log.Fatalf("unknown report name: %s", reportName)
		}
		return
	}
	log.Debugf("CaptureId %s==%s", captureId, capId)
	if capId != view.EmptyString {
		rdr := readers.NewCaptureReader(headersCache, packetCache, peersCache, pdb, sysInfo.GetWorkDir())
		if s3 == nil || !sysInfo.IsMinioStorageActive() {
			// override mode - no cloud storage
			err = rdr.ReadCaptureDir(capId, sysInfo.GetWorkDir())
			if err != nil {
				log.Errorf("unable to read capture %s from working directory %s. Error: %v", capId, sysInfo.GetWorkDir(), err)
			}
		} else {
			// override mode - use cloud storage
			log.Debugf("MAIN readers.ProcessCaptureFiles %s", capId)
			fileCount, err := s3.ProcessCaptureFiles(capId, rdr)
			if err != nil {
				log.Errorf("unable to process capture %s from cloud storage. Error: %v", capId, err)
			}
			log.Printf("%d file(s) procesed", fileCount)
		}
		return
	}
	log.Println("entering service mode")
	// service mode
	ws := controllers.NewService(entities.WebServiceConfig{
		APIkey:         sysInfo.GetAPIKey(),
		ProductionMode: sysInfo.IsProductionMode(),
		WorkDir:        sysInfo.GetWorkDir(),
		AgentName:      sysInfo.GetAgentName(),
	}, headersCache, packetCache, peersCache, s3, pdb, sysInfo.GetNamespace(), sysInfo.GetWorkspace(), apihubClient)
	r := mux.NewRouter()
	r.SkipClean(true)
	r.UseEncodedPath()
	// set API handlers
	r.HandleFunc(view.LoadPath, ws.OnCaptureLoad).Methods(http.MethodGet)
	r.HandleFunc(view.LoadStatusReportPath, ws.OnCaptureLoadStatus).Methods(http.MethodGet)
	r.HandleFunc(view.ServiceOperationsReportPath, ws.OnServiceOperationsReportGenerate).Methods(http.MethodPost) // generate
	r.HandleFunc(view.ServiceOperationsRenderPath, ws.OnServiceOperationsReportOutput).Methods(http.MethodGet)    // send it out
	r.HandleFunc(view.MinioDeleteCapturePath, ws.OnCaptureDelete).Methods(http.MethodDelete)                      // send it out
	if !sysInfo.IsProductionMode() {
		r.HandleFunc(view.MinioCleanupCapturePath, ws.OnCaptureCleanup).Methods(http.MethodDelete) // send it out
		r.PathPrefix("/debug/").Handler(http.DefaultServeMux)
		//r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		//r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		//r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		//r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}
	// set TTL reactions
	r.HandleFunc("/live", ws.OnStatus).Methods(http.MethodGet)
	r.HandleFunc("/ready", ws.OnStatus).Methods(http.MethodGet)
	r.HandleFunc("/startup", ws.OnStatus).Methods(http.MethodGet)
	srv := makeServer(sysInfo, r)
	log.Fatalf("Service fatal error:%v", srv.ListenAndServe())
}
