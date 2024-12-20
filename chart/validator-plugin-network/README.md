
Validator-plugin-network
===========

Perform various network validations (DNS, ICMP, MTU, Netcat, IP-range)


## Configuration

The following table lists the configurable parameters of the Validator-plugin-network chart and their default values.

| Parameter                | Description             | Default        |
| ------------------------ | ----------------------- | -------------- |
| `controllerManager.manager.args` |  | `["--health-probe-bind-address=:8081", "--metrics-bind-address=:8443", "--leader-elect"]` |
| `controllerManager.manager.containerSecurityContext.allowPrivilegeEscalation` |  | `true` |
| `controllerManager.manager.containerSecurityContext.capabilities.add` |  | `["NET_RAW"]` |
| `controllerManager.manager.containerSecurityContext.capabilities.drop` |  | `["ALL"]` |
| `controllerManager.manager.image.repository` |  | `"quay.io/validator-labs/validator-plugin-network"` |
| `controllerManager.manager.image.tag` | x-release-please-version | `"v0.1.2"` |
| `controllerManager.manager.resources.limits.cpu` |  | `"500m"` |
| `controllerManager.manager.resources.limits.memory` |  | `"128Mi"` |
| `controllerManager.manager.resources.requests.cpu` |  | `"10m"` |
| `controllerManager.manager.resources.requests.memory` |  | `"64Mi"` |
| `controllerManager.replicas` |  | `1` |
| `controllerManager.serviceAccount.annotations` |  | `{}` |
| `kubernetesClusterDomain` |  | `"cluster.local"` |
| `metricsService.ports` |  | `[{"name": "https", "port": 8443, "protocol": "TCP", "targetPort": 8443}]` |
| `metricsService.type` |  | `"ClusterIP"` |
| `env` |  | `[]` |
| `proxy.enabled` |  | `false` |
| `proxy.image` |  | `"quay.io/validator-labs/validator-certs-init:latest"` |
| `proxy.secretName` |  | `"proxy-cert"` |



---
_Documentation generated by [Frigate](https://frigate.readthedocs.io)._

