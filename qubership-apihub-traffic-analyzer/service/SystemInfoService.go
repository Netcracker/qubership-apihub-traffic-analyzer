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

package service

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Netcracker/qubership-apihub-traffic-analyzer/qubership-apihub-traffic-analyzer/view"
	log "github.com/sirupsen/logrus"
)

const (
	BasePath             = "BASE_PATH"
	ProductionMode       = "PRODUCTION_MODE"
	LogLevel             = "LOG_LEVEL"
	ListenAddress        = "LISTEN_ADDRESS"
	OriginAllowed        = "ORIGIN_ALLOWED"
	APIkey               = "TRAFFIC_API_KEY"
	ApiHubAgentName      = "API_HUB_AGENT"
	ApiHubUrl            = "APIHUB_URL"
	ApiHubAccessToken    = "APIHUB_ACCESS_TOKEN"
	PgHost               = "APIHUB_POSTGRESQL_HOST"
	PgPort               = "APIHUB_POSTGRESQL_PORT"
	PgDb                 = "APIHUB_TRAFFIC_POSTGRESQL_DB"
	PgUser               = "APIHUB_TRAFFIC_POSTGRESQL_USERNAME"
	PgPassword           = "APIHUB_TRAFFIC_POSTGRESQL_PASSWORD"
	PgSslMode            = "APIHUB_SNIFFER_POSTGRESQL_SSL_MODE"
	InsecureProxy        = "INSECURE_PROXY"
	MinioAccessKeyId     = "STORAGE_SERVER_ACCESS_KEY_ID"
	MinioSecretAccessKey = "STORAGE_SERVER_SECRET_ACCESS_KEY"
	MinioCrt             = "STORAGE_SERVER_CRT"
	MinioEndpoint        = "STORAGE_SERVER_ENDPOINT"
	MinioBucketName      = "STORAGE_SERVER_BUCKET_NAME"
	MinioStorageActive   = "MINIO_STORAGE_ACTIVE"
	WorkDir              = "WORK_DIR"
	CaptureId            = "CAPTURE_ID"
	SchemaName           = "PG_SCHEMA_NAME"
	KubeNamespace        = "NAMESPACE"
	WorkSpace            = "WORKSPACE"
	paramError           = "mandatory parameter %s is empty"
	defPgPort            = 5432
	defDotDir            = "."
	defLocalHost         = "localhost"
	DefListenaddress     = ":8080"
	DefApiHubAgentName   = "k8s-apps3_api-hub-dev" // "k8sApps3-api-hub-dev"
)

func NewSystemInfoService() (SystemInfoService, error) {
	s := &systemInfoServiceImpl{
		systemInfoMap: make(map[string]interface{})}
	if err := s.Init(); err != nil {
		log.Error("Failed to read system info: " + err.Error())
		return nil, err
	}
	return s, nil
}

type SystemInfoService interface {
	Init() error
	GetWorkDir() string
	//	GetJwtPrivateKey() []byte
	IsProductionMode() bool
	GetLogLevel() string
	GetListenAddress() string
	GetPGHost() string
	GetPGPort() int
	GetPGDB() string
	GetPGUser() string
	GetPGPassword() string
	GetPGSSLMode() string
	GetCredsFromEnv() *view.DbCredentials
	GetMinioAccessKeyId() string
	GetMinioSecretAccessKey() string
	GetMinioCrt() string
	GetMinioEndpoint() string
	GetMinioBucketName() string
	IsMinioStorageActive() bool
	GetMinioStorageCreds() *view.MinioStorageCreds
	IsMinioStoreOnlyBuildResult() bool
	GetCaptureId() string
	GetBasePath() string
	CmdLineOverride(captureId, baseDir, workDir string, connAttrs view.DbCredentials)
	Validated() error
	GetOriginAllowed() string
	GetAPIKey() string
	GetApiHubUrl() string
	GetApiHubAccessToken() string
	GetWorkspace() string
	GetNamespace() string
	GetAgentName() string
}
type systemInfoServiceImpl struct {
	systemInfoMap map[string]interface{}
}

func (g *systemInfoServiceImpl) GetCredsFromEnv() *view.DbCredentials {
	return &view.DbCredentials{
		Host:     g.GetPGHost(),
		Port:     g.GetPGPort(),
		Database: g.GetPGDB(),
		Username: g.GetPGUser(),
		Password: g.GetPGPassword(),
		SSLMode:  g.GetPGSSLMode(),
	}
}

func (g *systemInfoServiceImpl) GetPGHost() string {
	return g.getString(PgHost)
}

func (g *systemInfoServiceImpl) getString(varName string) string {
	v, found := g.systemInfoMap[varName]
	if found {
		return v.(string)
	}
	return view.EmptyString
}

func (g *systemInfoServiceImpl) GetPGPort() int {
	return g.getInt(PgPort, defPgPort)
}

func (g *systemInfoServiceImpl) getInt(varName string, defVal int) int {
	v, found := g.systemInfoMap[varName]
	if found {
		return v.(int)
	}
	return defVal

}

func (g *systemInfoServiceImpl) getBool(varName string, defVal bool) bool {
	v, found := g.systemInfoMap[varName]
	if found {
		return v.(bool)
	}
	return defVal

}

func (g *systemInfoServiceImpl) GetPGDB() string {
	return g.getString(PgDb)
}

func (g *systemInfoServiceImpl) GetPGUser() string {
	return g.getString(PgUser)
}

func (g *systemInfoServiceImpl) GetPGPassword() string {
	return g.getString(PgPassword)
}

func (g *systemInfoServiceImpl) GetPGSSLMode() string {
	return g.getString(PgSslMode)
}

func (g *systemInfoServiceImpl) GetMinioStorageCreds() *view.MinioStorageCreds {
	bucket := g.GetMinioBucketName()
	return &view.MinioStorageCreds{
		BucketName:      bucket,
		IsActive:        g.IsMinioStorageActive(),
		Endpoint:        g.GetMinioEndpoint(),
		Crt:             g.GetMinioCrt(),
		AccessKeyId:     g.GetMinioAccessKeyId(),
		SecretAccessKey: g.GetMinioSecretAccessKey(),
		WorkDir:         g.GetWorkDir(),
	}
}

func (g *systemInfoServiceImpl) GetMinioBucketName() string {
	return g.getString(MinioBucketName)
}

func (g *systemInfoServiceImpl) IsMinioStorageActive() bool {
	return g.getBool(MinioStorageActive, false)
}

func (g *systemInfoServiceImpl) GetMinioEndpoint() string {
	return g.getString(MinioEndpoint)
}

func (g *systemInfoServiceImpl) GetMinioCrt() string {
	return g.getString(MinioCrt)
}
func (g *systemInfoServiceImpl) GetMinioSecretAccessKey() string {
	return g.getString(MinioSecretAccessKey)
}
func (g *systemInfoServiceImpl) IsMinioStoreOnlyBuildResult() bool {
	return false
}

func (g *systemInfoServiceImpl) GetMinioAccessKeyId() string {
	return g.getString(MinioAccessKeyId)
}

// CmdLineOverride override some values from command line
func (g *systemInfoServiceImpl) CmdLineOverride(captureId, baseDir, workDir string, connAttrs view.DbCredentials) {
	g.fromString(CaptureId, captureId)
	g.fromString(WorkDir, workDir)
	g.fromString(BasePath, baseDir)
	g.fromString(PgHost, connAttrs.Host)
	g.fromString(PgUser, connAttrs.Username)
	g.fromString(PgPassword, connAttrs.Password)
	g.fromString(PgSslMode, connAttrs.SSLMode)
	g.fromString(PgDb, connAttrs.Database)
	g.fromString(SchemaName, connAttrs.Schema)
	if connAttrs.Port > 0 {
		g.systemInfoMap[PgPort] = connAttrs.Port
	}
}

func (g *systemInfoServiceImpl) fromString(name, value string) {
	if value != view.EmptyString {
		g.systemInfoMap[name] = value
	}
}

// fromEnv
// extracts string value from an environment variable
func (g *systemInfoServiceImpl) fromEnv(envName string, defVal string) {
	sVal := os.Getenv(envName)
	if sVal == view.EmptyString {
		sVal = defVal
	}
	g.systemInfoMap[envName] = sVal
}

// fromEnvInt
// extracts integer value from an environment variable
func (g *systemInfoServiceImpl) fromEnvInt(envName string, defVal int) {
	sVal := os.Getenv(envName)
	if sVal != view.EmptyString {
		i, e := strconv.ParseInt(sVal, 10, 64)
		if e == nil {
			defVal = int(i)
		} else {
			log.Errorf("non numeric value '%s' found in environment variable %s", sVal, envName)
		}
	}
	g.systemInfoMap[envName] = defVal
}

// fromEnvBool
// extracts boolean value from an environment variable
func (g *systemInfoServiceImpl) fromEnvBool(envName string, defVal bool) {
	sVal := os.Getenv(envName)
	if sVal != view.EmptyString {
		i, e := strconv.ParseBool(sVal)
		if e == nil {
			defVal = i
		} else {
			log.Errorf("non boolean value '%s' found in environment variable %s", sVal, envName)
		}
	}
	g.systemInfoMap[envName] = defVal
}

// Init
// read and interpret environment values
func (g *systemInfoServiceImpl) Init() error {
	// list of parameters
	strValues := []string{CaptureId, PgUser, PgPassword, PgSslMode, PgDb, SchemaName, ListenAddress, APIkey,
		LogLevel, OriginAllowed, MinioCrt, MinioAccessKeyId, MinioEndpoint, MinioBucketName, MinioSecretAccessKey,
		ApiHubAccessToken, ApiHubUrl, KubeNamespace, WorkSpace, ApiHubAgentName}
	// those will be initialized as empty strings
	for _, svn := range strValues {
		g.fromEnv(svn, view.EmptyString)
	}
	// strings with non-empty defaults
	if g.getString(ListenAddress) == view.EmptyString {
		log.Warnf("Listen address not found in environment variable %s. Defaulting to %s", ListenAddress, DefListenaddress)
		g.fromString(ListenAddress, DefListenaddress)
	}
	if g.getString(ApiHubAgentName) == view.EmptyString {
		g.fromString(ApiHubAgentName, DefApiHubAgentName)
	}
	g.fromEnv(WorkDir, os.TempDir())
	g.fromEnv(BasePath, defDotDir)
	g.fromEnv(PgHost, defLocalHost)
	g.fromEnv(PgSslMode, "off")
	// numeric
	g.fromEnvInt(PgPort, defPgPort)
	// booleans
	g.fromEnvBool(ProductionMode, true)
	g.fromEnvBool(InsecureProxy, false)
	g.fromEnvBool(MinioStorageActive, false)
	return nil
}

func (g *systemInfoServiceImpl) GetWorkDir() string {
	return g.getString(WorkDir)
}

func (g *systemInfoServiceImpl) GetLogLevel() string {
	return g.getString(LogLevel)
}

func (g *systemInfoServiceImpl) IsProductionMode() bool {
	return g.getBool(ProductionMode, false)
}

func (g *systemInfoServiceImpl) GetListenAddress() string {
	return g.getString(ListenAddress)
}

func (g *systemInfoServiceImpl) GetCaptureId() string {
	return g.getString(CaptureId)
}

func (g *systemInfoServiceImpl) GetBasePath() string {
	return g.getString(BasePath)
}

func (g *systemInfoServiceImpl) Validated() error {
	nonEmpty := []string{WorkDir, BasePath, PgHost, PgUser, PgPassword, PgDb, MinioEndpoint, MinioBucketName,
		KubeNamespace, WorkSpace, ApiHubUrl, ApiHubAccessToken}
	for _, constraintValue := range nonEmpty {
		if g.getString(constraintValue) == view.EmptyString {
			return fmt.Errorf(paramError, constraintValue)
		}
	}
	if g.getInt(PgPort, -1) < 0 {
		return fmt.Errorf(paramError, PgPort)
	}
	return nil
}

func (g *systemInfoServiceImpl) GetOriginAllowed() string {
	return g.getString(OriginAllowed)
}

func (g *systemInfoServiceImpl) GetAPIKey() string {
	return g.getString(APIkey)
}

func (g *systemInfoServiceImpl) GetApiHubUrl() string {
	return g.getString(ApiHubUrl)
}

func (g *systemInfoServiceImpl) GetApiHubAccessToken() string {
	return g.getString(ApiHubAccessToken)
}
func (g *systemInfoServiceImpl) GetWorkspace() string {
	return g.getString(WorkSpace)
}
func (g *systemInfoServiceImpl) GetNamespace() string {
	return g.getString(KubeNamespace)
}
func (g *systemInfoServiceImpl) GetAgentName() string {
	return g.getString(ApiHubAgentName)
}
