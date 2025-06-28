# API HUB Traffic analyser

## Installation notes

### Preinstalled software

The following software should be installed before APIHUB installation:

| Dependency  | Minimal version | Mandatory/Optional | Comments                                  |
|-------------|-----------------|--------------------|-------------------------------------------|
| PostgreSQL  | 13              | Mandatory          |                                           |
| Minio       | 1.2.3           | Optional           | For store cold data and reduce load to PG |

## HWE

### APIHUB traffic analyzer

|              | CPU request | CPU limit | RAM request | RAM limit |
|--------------|-------------|-----------|-------------|-----------|
| Dev profile  | 300m        | 1         | 500Mi       | 500Mi     |
| Prod profile | 600m        | 2         | 1024Mi      | 1024Mi    |

## Deploy parameters

### APIHUB traffic analyzer parameters description

A complete configuration value list is shown below:

| name                               | Default Value         | description                                                                                                            |
|------------------------------------|-----------------------|------------------------------------------------------------------------------------------------------------------------|
| BASE_PATH                          | .                     | A path to database migration scripts (SQL)                                                                             |
| PRODUCTION_MODE                    | false                 | A production mode indicator (usually restricts logging and diagnostic functions)                                       |
| LOG_LEVEL                          | info                  | increase or decrease logging messages counts. one of: trace, debug, info, warning, error, fatal, panic                 |
| LISTEN_ADDRESS                     | :8080                 | Interface endpoint address                                                                                             |
| API_HUB_AGENT                      | k8s-apps3_api-hub-dev | A name to query APIHUB service (package) operations report                                                             |
| APIHUB_URL                         |                       | An address where APIHUB backend is listening on                                                                        |
| APIHUB_ACCESS_TOKEN                |                       | An API key to query APIHUB backend                                                                                     |
| APIHUB_POSTGRESQL_HOST             |                       | PostgreSQL server host name or IP address                                                                              |
| APIHUB_POSTGRESQL_PORT             | 5432                  | PostgreSQL server port number                                                                                          |
| APIHUB_TRAFFIC_POSTGRESQL_DB       | apihub                | PostgreSQL database name (will be generated automatically when not set)                                                |
| APIHUB_TRAFFIC_POSTGRESQL_USERNAME | apihub                | PostgreSQL user name (will be generated automatically when not set)                                                    |
| APIHUB_TRAFFIC_POSTGRESQL_PASSWORD | APIhub1234            | PostgreSQL user password  (will be generated automatically when not set)                                               |
| APIHUB_POSTGRESQL_SSL_MODE         | off                   | PostgreSQL SSL mode (__*off*__ for the most cases)                                                                     |
| INSECURE_PROXY                     | false                 | Set to true to enable apihub playground work without authorization.                                                    |
| ORIGIN_ALLOWED                     |                       |                                                                                                                        |
| STORAGE_SERVER_USERNAME            |                       | Minio/S3 access key Id                                                                                                 |
| STORAGE_SERVER_PASSWORD            |                       | Minio/S3 access key                                                                                                    |
| STORAGE_SERVER_CRT                 |                       | Minio/S3 certificate                                                                                                   |
| STORAGE_SERVER_URL                 |                       | Minio/S3 service endpoint address                                                                                      |
| STORAGE_SERVER_BUCKET_NAME         |                       | Minio/S3 bucket name                                                                                                   |
| MINIO_STORAGE_ACTIVE               | false                 | Minio/S3 interface state (true for active). if PRODUCTION_MODE is set to __*true*__ then this value must be __*true*__ | 
| WORK_DIR                           | .                     | A local filesystem path used to create intermediate files and download the data from Minio/S3                          |
| NAMESPACE                          |                       | A namespace where POD is running                                                                                       |
| WORKSPACE                          |                       | A workpace name to query service (package) operations report                                                           |

## Command line to override parameters

Some configuration parameters which could be overridden or set with the command line options listed in the table below:

| command line option | parameter name                     | description                                                                                                                                         |
|---------------------|------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| -base-dir           | BASE_PATH                          | A path to database migration scripts (SQL)                                                                                                          |
| -capture-id         |                                    | To set Capture Id manually for testing purpose. If the value set then traffic analyser will read capture from local filesystem or S3/Minio and exit |
| -host               | APIHUB_POSTGRESQL_HOST             | PostgreSQL server host name or IP address                                                                                                           |
| -port               | APIHUB_POSTGRESQL_PORT             | PostgreSQL port number                                                                                                                              |
| -instance           | APIHUB_TRAFFIC_POSTGRESQL_DB       | PostgreSQL database                                                                                                                                 |
| -schema             |                                    | PostgreSQL schema to maintain more than one data set at single database server                                                                      | 
| -user               | APIHUB_TRAFFIC_POSTGRESQL_USERNAME | PostgreSQL user name                                                                                                                                |
| -password           | APIHUB_TRAFFIC_POSTGRESQL_PASSWORD | PostgreSQL user password                                                                                                                            |
| -ssl-mode           | APIHUB_POSTGRESQL_SSL_MODE         | PostgreSQL SSL mode (__*off*__ for the most cases)                                                                                                  |
| -work-dir           | WORK_DIR                           | A local filesystem path used to create intermediate files and download the data from Minio/S3                                                       |
| -report-name        |                                    | An URL of the report to receive                                                                                                                     |
| -service-name       |                                    | A service name for report (report parameter)                                                                                                        |
| -service-version    |                                    | A service version for report (report parameter). Makes sense with the service name only                                                             |
| -log-level          | info                               | A logging level: (trace, debug, info, warning, error, fatal, panic)                                                                                 |
