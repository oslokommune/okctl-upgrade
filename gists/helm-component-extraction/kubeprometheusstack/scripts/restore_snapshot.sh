#!/usr/bin/env bash

## Usage
# ./restore_snapshot.sh <snapshot_directory_name>

SNAPSHOT_DIRECTORY_NAME=$1
NS="monitoring"

DRY_RUN=client

pausePrometheus() {
    echo "💀 Turning off the Prometheus instance with replicas: 0"
    kubectl patch --namespace=${NS} --dry-run=${DRY_RUN} Prometheus kube-prometheus-stack-prometheus --type="json" \
        --patch '[{"op": "replace", "path": "/spec/replicas", value: 0}]'
}

resumePrometheus() {
    echo "🔥 Turning on the Prometheus instance with replicas: 1"
    kubectl patch --namespace=${NS} --dry-run=${DRY_RUN} Prometheus kube-prometheus-stack-prometheus --type="json" \
        --patch '[{"op": "replace", "path": "/spec/replicas", value: 1}]'
}

createDebuggerPod() {
    echo "🤖 Creating debugger pod"
    kubectl apply --namespace=${NS} --dry-run=${DRY_RUN} -f templates/debugger.yaml

    echo "🕐 Waiting for debugger pod to be ready"
    if [[ ${DRY_RUN} == "none" ]]; then
        kubectl wait --namespace=${NS} --for=condition=Ready pod -l app=debugger --timeout=180s # copiloted, verify
    else
        echo "🕐 Pretending to be waiting.."
    fi
}

deleteDebuggerPod() {
    echo "🧹 Cleaning up debugger pod"
    kubectl delete --namespace=${NS} --dry-run=${DRY_RUN} -f templates/debugger.yaml
}

cleanOutPrometheusData() {
    echo "🧹 Cleaning out Prometheus data"
    if [[ ${DRY_RUN} == "none" ]]; then
        kubectl exec --namespace=${NS} -it debugger -- rm -rf /prometheus/prometheus-db
    else
        echo "🧹 Pretending to be cleaning.."
    fi
}

restoreSnapshot() {
    echo "📦 Copying snapshot to Prometheus PVC"
    if [[ ${DRY_RUN} == "none" ]]; then
        kubectl cp --namespace=${NS} ${SNAPSHOT_DIRECTORY_NAME} debugger:/prometheus/prometheus-db
    else
        echo "📦 Pretending to be copying.."
    fi
}

pausePrometheus
createDebuggerPod
cleanOutPrometheusData
restoreSnapshot
deleteDebuggerPod
resumePrometheus
