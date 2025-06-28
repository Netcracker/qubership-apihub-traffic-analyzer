# Qubership APIHUB Traffic Analyzer Helm Chart

## Prerequisites

1. kubectl installed and configured for k8s cluster access. Namespace admin permissions required.
1. Helm installed
1. Supported k8s version - 1.23+

# 3rd party dependencies

| Name | Version | Mandatory/Optional | Comment |
| ---- | ------- |------------------- | ------- |
| S3 | Any | Mandatory | For reading captured data |
| Qubership APIHUB | Any | Mandatory | For reading API specifications data which is required for reports generation |
| Postgres | 14+ | Mandatory | For store analysis results data |

## HWE

|     | CPU request | CPU limit | RAM request | RAM limit |
| --- | ----------- | --------- | ----------- | --------- |
| Minimal level        | 30m | 300m   | 256Mi | 256Mi |
| Average load level   | 30m | 2      | 2Gi   | 4Gi   |


## Set up values.yml

1. Download Qubership APIHUB Traffic Analyzer helm chart
1. Fill `values.yaml` with corresponding deploy parameters. `values.yaml` is self-documented, so please refer to it

## Execute helm install

In order to deploy Qubership APIHUB to your k8s cluster execute the following command:

```
helm install qubership-apihub-traffic-analyzer -n qubership-apihub-traffic-analyzer --create-namespace -f ./helm-templates/qubership-apihub-traffic-analyzer/values.yaml ./helm-templates/qubership-apihub-traffic-analyzer
```

In order to uninstall Qubership APIHUB Traffic Analyzer from your k8s cluster execute the following command:

```
helm uninstall qubership-apihub-traffic-analyzer -n qubership-apihub-traffic-analyzer
```