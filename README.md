Forked for Polygon.io, <a href="https://polygon.io">Get your free Stock Data API Key</a>

# drone-plugin-kube

[![](https://images.microbadger.com/badges/version/danielgormly/drone-plugin-kube.svg)](https://microbadger.com/images/danielgormly/drone-plugin-kube "Get your own version badge on microbadger.com")

A simple Drone plugin for updating Kubernetes Deployments from templates & ConfigMaps from files.
Follows from [vallard/drone-kube](https://github.com/vallard/drone-kube) and then [drone-plugin-kube](https://github.com/danielgormly/drone-plugin-kube)
but adds support for services and ingresses. It also has a slightly better working deployment watcher.

Usage: See [drone.md](./drone.md)

## Deployment templates

Deployment config files are first interpreted by **aymerick/raymond** ([handlebarsjs](http://handlebarsjs.com/) equivalent). You can use all available raymond expressions and anything you put in settings will be made available in your deployment template e.g. `{{namespace}}`. See [example/deployment.template.yaml](/example/deployment.template.yaml) for a complete example.

## Config maps from files

In this case, you can create a template just like deployment.yaml but you can provide a file path (relative to the repo's root) in the plugin setting `configmap_file`. (Currently only accepts utf-8 encoded data). Like deployments, this will both create new or update existing configmaps (based on the configmap name).

#### Adding a service account to Kubernetes that can manage deployments
See [example/Role.yaml](example/Role.yaml), [example/ServiceAccount.yaml](example/ServiceAccount.yaml), [example/RoleBinding.yaml](example/RoleBinding.yaml).

## Development notes
- Kubernetes client not yet supported by dep, so we are using
[`brew install glide`](https://github.com/Masterminds/glide).
- Update dependencies with brew `glide update --strip-vendor`
- [kubernetes/client-go installation notes](https://github.com/kubernetes/client-go/blob/master/INSTALL.md)
- [Creating a Drone plugin in Go](https://docs.drone.io/plugins/tutorials/golang/)
- [Client-go API Docs @ godoc.org](https://godoc.org/k8s.io/client-go/kubernetes)
- Testing with minikube (OSX: `brew cask install minikube`)
