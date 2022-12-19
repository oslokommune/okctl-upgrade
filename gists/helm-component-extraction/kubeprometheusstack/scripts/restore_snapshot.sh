#!/usr/bin/env bash

## Usage
# ./restore_snapshot.sh <snapshot_directory_name>

SNAPSHOT_DIRECTORY_NAME=$1
ns="monitoring"

## Turn off the Prometheus instance
kubectl ${ns} --namespace patch Prometheus kube-prometheus-stack-prometheus \
	--patch '[{"op": "replace", "path": "/spec/replicas", value: 0}]'

## Create debugger pod
kubectl --namespace ${ns} apply -f ../templates/debugger.yaml
kubectl --namespace ${ns} wait --for=condition=Ready pod -l app=debugger --timeout=60s # copiloted, verify

## Copy snapshot to Prometheus PVC
kubectl --namespace ${ns} cp ${SNAPSHOT_DIRECTORY_NAME} debugger:/prometheus

## Place snapshot in correct location
kubectl --namespace ${ns} exec -it debugger -- rm -rf /prometheus/prometheus-db/*
kubectl --namespace ${ns} exec -it debugger -- mv -r /prometheus/${SNAPSHOT_DIRECTORY_NAME}/* /prometheus/prometheus-db

## Turn on the Prometheus instance
kubectl ${ns} --namespace patch Prometheus kube-prometheus-stack-prometheus \
	--patch '[{"op": "replace", "path": "/spec/replicas", value: 1}]'
