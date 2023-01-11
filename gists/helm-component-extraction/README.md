# Ejecting components from your okctl cluster

## Motivation

[Slack post](https://oslokommune.slack.com/archives/C02U8HJARLH/p1668410829488739)

## Recommended extraction order

1. [Kube-Prometheus-Stack](./kubeprometheusstack/) // Owns the ServiceMonitor resource other components rely upon
2. [Promtail](./promtail/)
3. [Loki](./loki/) // Depends on Promtail
4. [ArgoCD](./argocd/)
2. [Autoscaler](./autoscaler/)
6. AWS Loadbalancer Controller (TBA)
7. External DNS (TBA)
