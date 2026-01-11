# Kubernetes

Helm chart and Kubernetes deployment assets for Conductor.

## Planned chart features
- Deployment/StatefulSet for Conductor backend.
- Optional ingress + service exposure.
- ConfigMaps/Secrets for configuration and credentials.
- Leader election via Lease API.
- Prometheus ServiceMonitor support.

## Security & RBAC
- Kubernetes RBAC for required resources.
- App-level RBAC (admin/dba/monitor/viewer) handled in the service layer.
- Optional OAuth2/OIDC integration for authentication.
