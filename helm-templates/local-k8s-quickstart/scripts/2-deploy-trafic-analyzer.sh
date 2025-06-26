echo "---Start TRAFFIC-ANALYZER deploy using Helm---"
helm install qubership-apihub-traffic-analyzer -n qubership-apihub-traffic-analyzer --create-namespace -f ../qubership-apihub-traffic-analyzer/local-k8s-values.yaml -f ../qubership-apihub-traffic-analyzer/local-secrets.yaml ../../qubership-apihub-traffic-analyzer
echo "---Complete TRAFFIC-ANALYZER deploy using Helm---"