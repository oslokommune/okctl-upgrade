#!/usr/bin/env bash

## Usage
# ./restore_snapshot.sh <snapshot_directory_name>

SNAPSHOT_DIRECTORY_NAME=$1
NS="monitoring"

DRY_RUN=client

## Turn off the Prometheus instance
echo "💀 Turning off the Prometheus instance with replicas: 0"
kubectl patch --namespace=${NS} --dry-run=${DRY_RUN} Prometheus kube-prometheus-stack-prometheus --type="json" \
	--patch '[{"op": "replace", "path": "/spec/replicas", value: 0}]'

## Create debugger pod
echo "🤖 Creating debugger pod"
kubectl apply --namespace ${NS} --dry-run=${DRY_RUN} -f templates/debugger.yaml
echo "🕐 Waiting for debugger pod to be ready"
if [[ ${DRY_RUN} == "none" ]]; then
	kubectl wait --namespace ${NS} --for=condition=Ready pod -l app=debugger --timeout=60s # copiloted, verify
else
	echo "🕐 Waiting.."
fi

## Copy snapshot to Prometheus PVC
echo "📦 Copying snapshot to Prometheus PVC"
if [[ ${DRY_RUN} == "none" ]]; then
	kubectl cp --namespace ${NS} ${SNAPSHOT_DIRECTORY_NAME} debugger:/prometheus
else
	echo "📦 Copying.."
fi

## Place snapshot in correct location
echo "🧹 Cleaning out generated Prometheus data"
if [[ ${DRY_RUN} == "none" ]]; then
	kubectl exec --namespace ${NS} -it debugger -- rm -rf /prometheus/prometheus-db/*
else
	echo "🧹 Cleaning.."
fi

echo "📂 Placing snapshot in correct location"
if [[ ${DRY_RUN} == "none" ]]; then
	kubectl exec --namespace ${NS} -it debugger -- mv /prometheus/${SNAPSHOT_DIRECTORY_NAME}/* /prometheus/prometheus-db
else
	echo "📂 Moving.."
fi

## Turn on the Prometheus instance
echo "🔥 Turning on the Prometheus instance with replicas: 1"
kubectl patch --namespace ${NS} --dry-run=${DRY_RUN} Prometheus kube-prometheus-stack-prometheus --type="json" \
	--patch '[{"op": "replace", "path": "/spec/replicas", value: 1}]'
