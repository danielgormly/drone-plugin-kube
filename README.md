# drone-plugin-kube

[![](https://images.microbadger.com/badges/version/danielgormly/drone-plugin-kube.svg)](https://microbadger.com/images/danielgormly/drone-plugin-kube "Get your own version badge on microbadger.com")

A simple Drone plugin for updating Kubernetes resources from yaml templates. Follows from [vallard/drone-kube](https://github.com/vallard/drone-kube). The plugin will create a resource if it doesn't already exist i.e. behaves similar to `kubectl apply`. This plugin supports deployments, configmaps, ingresses, and services.

Usage: See [drone.md](./drone.md)

#### Creating a Kubernetes service account for Drone
See [example/Role.yaml](example/Role.yaml), [example/ServiceAccount.yaml](example/ServiceAccount.yaml), [example/RoleBinding.yaml](example/RoleBinding.yaml).

## Usage notes:

- The watching process is not currently reliable i.e. it doesn't properly wait for new deployments to become live. Not entirely sure how this should behave but I think behind a flag would make sense. PRs welcome.

## Development notes
- [kubernetes/client-go installation notes](https://github.com/kubernetes/client-go/blob/master/INSTALL.md)
- [Creating a Drone plugin in Go](https://docs.drone.io/plugins/golang/)
- [Client-go API Docs @ godoc.org](https://pkg.go.dev/k8s.io/client-go/kubernetes?tab=doc)
- Testing with minikube (OSX: `brew cask install minikube`)

## Acknowledgements
- [@vallard](https://github.com/vallard) for the original plugin
- [@jbonzo](https://github.com/jbonzo) for ingress, service support, improved error handling etc
