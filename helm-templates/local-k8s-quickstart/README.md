# Qubership APIHUB Traffic Analyzer utilities for quick deployment to local k8s cluster

For local k8s cluster set up please Please refer to [Qubership APIHUB quickstart Installation on local k8s cluster](https://github.com/Netcracker/qubership-apihub/tree/main/helm-templates/local-k8s-quickstart)

## Prerequisites

1. Helm
2. Bash (GitBash, Cygwin, etc)

Assumptions:

1. local k8s cluter is set up, up and running, kubectl configured
2. Postgres and APIHUB already installed to your k8s cluster. Refer to [corresponding guide](https://github.com/Netcracker/qubership-apihub/tree/main/helm-templates/local-k8s-quickstart)
3. Minio and Sniffer Agent already installed to your k8s cluster. Refer to [corresponding guide](https://github.com/Netcracker/qubership-apihub-sniffer-agent/tree/helm/helm-templates/local-k8s-quickstart)

## Deployment

Deployment phases represented by scripts in `scripts` folder, so you can see what happens on each step in them:

- `1-get-secrets.sh` - reads Minio and APIHUB access keys (traffic analyzer uses the same Minio instance and bucket as sniffer agent) for existing k8s deployments
- `2-deploy-trafic-analyzer.sh` - deploy Traffic Analyzer itself. This script includes Traffic Analyzer secrets generation.

One-liners:

1. `quickstart.sh` - Traffic Analyzer.

## Uninstallation

```
helm uninstall qubership-apihub-traffic-analyzer -n qubership-apihub-traffic-analyzer
```