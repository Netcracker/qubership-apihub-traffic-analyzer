echo "---Get secrets---"
export SNIFFER_API_KEY=$(kubectl get secrets  -n qubership-apihub-sniffer-agent qubership-apihub-sniffer-agent-keys-secret  -o jsonpath={.data.api_key})
export ACCESS_TOKEN=$(kubectl get secrets  -n qubership-apihub qubership-apihub-backend-access-token-secret -o jsonpath={.data.token} | base64 -d)

export STORAGE_SERVER_USERNAME=$(kubectl get secrets  -n qubership-apihub-sniffer-agent qubership-apihub-sniffer-agent-s3-secret  -o jsonpath={.data.storage_server_username} | base64 -d)
export STORAGE_SERVER_PASSWORD=$(kubectl get secrets  -n qubership-apihub-sniffer-agent qubership-apihub-sniffer-agent-s3-secret  -o jsonpath={.data.storage_server_password} | base64 -d)
export STORAGE_SERVER_CRT=$(kubectl get secrets  -n qubership-apihub-sniffer-agent qubership-apihub-sniffer-agent-s3-secret  -o jsonpath={.data.storage_server_crt} | base64 -d)
export STORAGE_SERVER_URL=$(kubectl get secrets  -n qubership-apihub-sniffer-agent qubership-apihub-sniffer-agent-s3-secret  -o jsonpath={.data.storage_server_url} | base64 -d)
export STORAGE_SERVER_BUCKETNAME=$(kubectl get secrets  -n qubership-apihub-sniffer-agent qubership-apihub-sniffer-agent-s3-secret  -o jsonpath={.data.storage_server_bucketname} | base64 -d)

envsubst < ../qubership-apihub-traffic-analyzer/local-secrets.yaml.template > ../qubership-apihub-traffic-analyzer/local-secrets.yaml 

echo "---APIHUB and SNIFFER AGENT integration secrets---"
echo "SNIFFER_API_KEY: $SNIFFER_API_KEY"
echo "APIHUB ACCESS_TOKEN: $ACCESS_TOKEN"
echo "---S3 integration secrets---"
echo "STORAGE_SERVER_BUCKETNAME: $STORAGE_SERVER_BUCKETNAME"
echo "STORAGE_SERVER_URL: $STORAGE_SERVER_URL"
echo "STORAGE_SERVER_USERNAME: $STORAGE_SERVER_USERNAME"
echo "STORAGE_SERVER_PASSWORD: $STORAGE_SERVER_PASSWORD"
echo "STORAGE_SERVER_CRT: $STORAGE_SERVER_CRT"


