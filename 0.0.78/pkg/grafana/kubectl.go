package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

const (
	grafanaRepository                = "grafana/grafana"
	expectedGrafanaVersionPreUpgrade = "7.3.5"
	upgradeTag                       = "7.5.12"
	monitoringNamespace              = "monitoring"
	grafanaDeploymentName            = "kube-prometheus-stack-grafana"
	grafanaContainerName             = "grafana"
)

func acquireKubectlClient(kubeConfigPath string) (*kubernetes.Clientset, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("creating rest config: %w", err)
	}

	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("initializing client: %w", err)
	}

	return client, nil
}

func getCurrentGrafanaVersion(clientSet *kubernetes.Clientset) (string, error) {
	grafanaContainerIndex, err := getContainerIndexByName(clientSet, grafanaContainerName)
	if err != nil {
		return "", fmt.Errorf("getting Grafana container index: %w", err)
	}

	result, err := clientSet.AppsV1().Deployments(monitoringNamespace).Get(
		context.Background(),
		grafanaDeploymentName,
		metav1.GetOptions{},
	)
	if err != nil {
		return "", fmt.Errorf("getting deployment: %w", err)
	}

	parts := strings.Split(result.Spec.Template.Spec.Containers[grafanaContainerIndex].Image, ":")

	return parts[1], nil
}

func patchGrafanaDeployment(clientSet *kubernetes.Clientset, dryRun bool) (Receipts, error) {
	dryRunOpts := []string{metav1.DryRunAll}
	if !dryRun {
		dryRunOpts = nil
	}

	receipts := NewReceiptsStack()

	receipts.Push("identifying relevant container")

	grafanaContainerIndex, err := getContainerIndexByName(clientSet, grafanaContainerName)
	if err != nil {
		return receipts, fmt.Errorf("acquiring Grafana container index: %w", err)
	}

	receipts.Push("generating upgrade patch")

	patch := Patch{
		Op:    jsonPatchOperationReplace,
		Path:  fmt.Sprintf("/spec/template/spec/containers/%d/image", grafanaContainerIndex),
		Value: fmt.Sprintf("%s:%s", grafanaRepository, upgradeTag),
	}

	raw, err := json.Marshal(patch)
	if err != nil {
		return receipts, fmt.Errorf("marshalling patch: %w", err)
	}

	receipts.Push("applying patch")

	_, err = clientSet.AppsV1().Deployments(monitoringNamespace).Patch(
		context.Background(),
		grafanaDeploymentName,
		types.JSONPatchType,
		raw,
		metav1.PatchOptions{DryRun: dryRunOpts},
	)
	if err != nil {
		return receipts, fmt.Errorf("patching Grafana deployment: %w", err)
	}

	receipts.Push("verifying new Grafana version")

	newVersion, err := getCurrentGrafanaVersion(clientSet)
	if err != nil {
		return receipts, fmt.Errorf("acquiring updated Grafana version: %w", err)
	}

	expectedVersion := upgradeTag
	if dryRun {
		expectedVersion = expectedGrafanaVersionPreUpgrade
	}

	err = validateVersion(expectedVersion, newVersion)
	if err != nil {
		return receipts, fmt.Errorf("validating new version: %w", err)
	}

	return receipts, nil
}

func getContainerIndexByName(clientSet *kubernetes.Clientset, name string) (int, error) {
	result, err := clientSet.AppsV1().Deployments(monitoringNamespace).Get(
		context.Background(),
		grafanaDeploymentName,
		metav1.GetOptions{},
	)
	if err != nil {
		return -1, fmt.Errorf("getting deployment data: %w", err)
	}

	for index, container := range result.Spec.Template.Spec.Containers {
		if container.Name == name {
			return index, nil
		}
	}

	return -1, fmt.Errorf("not found")
}
