# apihub-traffic-analyzer usage scenario

## Prerequisites

This scenario starts at the end of successful packet capture finish. Please consult apihub-sniffer-agent usage scenario.
All the available endpoints are described in [interface section](api/interface.yaml).

## Usage cycle

* Load and aggregate finished capture
* Delete raw capture data from S3 (non-mandatory step)
* Generate report or reports on capture data
* Receive (render) generated report data
* [TODO/concept] Clear aggregated report data

### Load and aggregate finished capture

Once the packet capturing was completed successfully use endpoint ```/api/v1/admin/capture/{captureId}/load``` for loading and aggregating raw capture data.
Pass capture id (a unique id that was provided by apihub-sniffer-agent) to start load and aggregation. The execution time depends on the data size linearly. 
Use interface ```/api/v1/admin/capture/{captureId}/status``` to receive data load status. Wait for the loading to complete before start generating reports. 

### Delete raw capture data from S3

Use endpoint ```/api/v1/admin/capture/{captureId}/delete``` to delete raw capture data that no longer required. Usually the operation finished quickly and removes S3/Minio objects related to the capture id, passed as a parameter.  

### Generate report

Use one of the endpoints ```/api/v1/report/*/generate``` to perform aggregations and calculations specific to the report.
Provide parameters (specific for each report type) with JSON in body, for example:

* capture id
* service name
* service version

The endpoint starts report generation in the background and returns unique report id immediately.
This unique report id will be required to receive/render report.

The results will be stored in the database. Different report data for different parameters can be stored in the database simultaneously. 

### Receive/render generated report data

Use one of the endpoints ```/api/v1/report/*/render``` to receive a report render. This render of the completed report will be created in different output formats (implemented for each report type separately):

* Microsoft Excel (.xlsx)
* JSON
* HTML
* XML