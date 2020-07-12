# sidecar-operator
Demo kubernetes operator for injecting sidecars

## Install Kubernetes operator framework SDK

See operator SDK installation:

https://sdk.operatorframework.io/docs/install-operator-sdk/

If using brew on MacOS run:

```
$ brew install operator-sdk
```

Install kustomize
```
$ brew install kustomize
```

## Bootstrapping the operator

To initialise a new operator run the command:

```
$ operator-sdk init --domain=bvwells.github.com --repo=github.com/bvwells/sidecar-operator
```

Adding a new custom resource definition and controller

```
$ operator-sdk create api --version=v1alpha1 --kind=SidecarOperator --group=sidecar 
```

Modify types api/v1alpha1/sidecaroperator_types.go and to update generated code run:

$ make generate

To update the CRD manifests run:

$ make manifests

## Deploy the CRD

Build the operator with and apply the CRD to you cluster run:

```
$ make install
```

To observe the CRD run:

```
$ kubectl get crds
NAME                                          CREATED AT
sidecaroperators.sidecar.bvwells.github.com   2020-07-12T20:30:25Z
```

## Run the operator locally

```
$ make run ENABLE_WEBHOOKS=false
```

## Deploy an example

Deploy an example CRD
```
$ kubectl apply -f <(echo "
apiVersion: sidecar.bvwells.github.com/v1alpha1
kind: SidecarOperator
metadata:
  name: sidecaroperator-sample
spec:
  image: "alpine:latest"
")
```

Check on the CRD by running:
```
$ kubectl get SidecarOperator
```

Get details on the created pod:
```
$ kubectl get pods
```

```
$ kubectl decribe pods example-sidecaroperator-pod
```

Get logs from running pod:
```
$ kubectl logs  example-sidecaroperator-pod 
```

## Deploy the Operator

```
$ make docker-build IMG=repository/bvwells/sidecar-operator:v0.0.1
```

Push the image to repository of choice.

```
$ make docker-push IMG=repository/bvwells/sidecar-operator:v0.0.1
```

## Deploy the operator

Set the namespace to run the operator. For example to run in the "default"
namespace run:

```
$ cd config/default/ && kustomize edit set namespace "default" && cd ../..
```

```
$ make deploy IMG=repository/bvwells/sidecar-operator:v0.0.1
```

## Prerequisites

The following tools are required to develop and test this operator example:
- git
- go version v1.12+.
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
