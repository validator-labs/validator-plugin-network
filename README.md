[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/validator-labs/validator-plugin-network/issues)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![Test](https://github.com/validator-labs/validator-plugin-network/actions/workflows/test.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/validator-labs/validator-plugin-network)](https://goreportcard.com/report/github.com/validator-labs/validator-plugin-network)
[![codecov](https://codecov.io/github/validator-labs/validator-plugin-network/graph/badge.svg?token=Q15XUCRNCN)](https://codecov.io/github/validator-labs/validator-plugin-network)
[![Go Reference](https://pkg.go.dev/badge/github.com/validator-labs/validator-plugin-network.svg)](https://pkg.go.dev/github.com/validator-labs/validator-plugin-network)

# validator-plugin-network
The Network [validator](https://github.com/validator-labs/validator) plugin ensures that your network matches a user-configurable expected state.

## Description

The Network validator plugin reconciles `NetworkValidator` custom resources to perform the following validations against your network:

1. Execute DNS lookups
2. Execute ICMP pings
3. Validate TCP connections to arbitrary host + port(s), optionally through an HTTP proxy
4. Check each IP in a given range (starting IP + next N IPs) to ensure that they're all unallocated
5. Check that the default NIC has an MTU of at least X, where X is the provided MTU
6. Check that each file in a list of URLs is available via an HTTP HEAD request, optionally with HTTP basic authentication.

Each `NetworkValidator` CR is (re)-processed every two minutes to continuously ensure that your network matches the expected state.

See the [samples](https://github.com/validator-labs/validator-plugin-network/tree/main/config/samples) directory for example `NetworkValidator` configurations.

## Adding CA certs

For HTTP file rules, you can add client CA certs to use by doing either of the following:

1. While installing the plugin, in `values.yaml`, you can set `proxy` to add a CA cert to the system cert pool.
2. While applying `NetworkValidator`s, in their specs, you can set `caCerts` to provide additional CA certs to be applied on top of the system cert pool. The certs can be provided inline or via secrets.

For TCP connection rules, you can add a client CA cert to use by doing the following:

1. While installing the plugin, in `values.yaml`, you can set `proxy` to add a CA cert to the system cert pool.

## Installation
The Network validator plugin is meant to be [installed by validator](https://github.com/validator-labs/validator/tree/gh_pages#installation) (via a ValidatorConfig), but it can also be installed directly as follows:

```bash
helm repo add validator-plugin-network https://validator-labs.github.io/validator-plugin-network
helm repo update
helm install validator-plugin-network validator-plugin-network/validator-plugin-network -n validator-plugin-network --create-namespace
```

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [kind](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/validator-plugin-network:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/validator-plugin-network:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
All contributions are welcome! Feel free to reach out on the [Spectro Cloud community Slack](https://spectrocloudcommunity.slack.com/join/shared_invite/zt-g8gfzrhf-cKavsGD_myOh30K24pImLA#/shared-invite/email).

Make sure `pre-commit` is [installed](https://pre-commit.com#install).

Install the `pre-commit` scripts:

```console
pre-commit install --hook-type commit-msg
pre-commit install --hook-type pre-commit
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

