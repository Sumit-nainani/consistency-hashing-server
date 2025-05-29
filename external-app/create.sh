#!/bin/bash

# Enable exit on error
set -e

# Define cluster name and config file with defaults
CLUSTER_NAME=${1:-my-cluster}
CONFIG_FILE=${2:-kind-config.yaml}

# Check if 'kind' is installed
if ! command -v kind &> /dev/null; then
    echo "❌ 'kind' is not installed. Please install kind first."
    exit 1
fi

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Config file '$CONFIG_FILE' not found!"
    exit 1
fi

# Check if a cluster with the same name already exists
if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
    echo "❌ A cluster named '$CLUSTER_NAME' already exists!"
    echo "➡️  You can delete it first with: kind delete cluster --name $CLUSTER_NAME"
    exit 1
fi

# Create the cluster
echo "🚀 Creating cluster '$CLUSTER_NAME' using config file '$CONFIG_FILE'..."
kind create cluster --name "$CLUSTER_NAME" --config "$CONFIG_FILE"

echo "✅ Cluster '$CLUSTER_NAME' created successfully!"

# Apply the namespaces YAML file
kubectl apply -f ../namespaces
echo "✅ Namespaces created: demo and monitoring."

# Apply Prometheus CRDs, operator, and Prometheus itself
kubectl apply -f ../prometheus-operator-crd
kubectl apply -f ../prometheus-operator
kubectl apply -f ../prometheus
echo "✅ Prometheus setup applied successfully."

kubectl apply -f ../kubernetes
echo "app deployment successfull."

kubectl apply -f ../prometheus-adapter/0-adapter
kubectl apply -f ../prometheus-adapter/1-custom-metrics
kubectl apply -f ../prometheus-adapter/2-resource-metrics
echo "prometheus adapter created."

echo "🚀 Installing metrics-server..."
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Check if metrics-server is installed successfully
echo "✅ Metrics server installed successfully."

kubectl patch deployment metrics-server -n kube-system --type='json' -p='[
  {"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"},
  {"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-preferred-address-types=InternalIP"}
]'
echo "✅ Patched metrics-server deployment with insecure TLS and preferred address types."

kubectl run curlpod --image=curlimages/curl --namespace=curl-pod --restart=Never -- sleep infinity
echo "✅ Curl pod created in 'curl-pod' namespace."

# Final confirmation
echo "✅ Cluster '$CLUSTER_NAME' setup complete with namespaces, Prometheus, and Metrics Server."