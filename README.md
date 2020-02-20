# drone-plugin-kube

[![](https://images.microbadger.com/badges/version/danielgormly/drone-plugin-kube.svg)](https://microbadger.com/images/danielgormly/drone-plugin-kube "Get your own version badge on microbadger.com")

A simple Drone plugin for updating Kubernetes resources from yaml templates. Follows from [vallard/drone-kube](https://github.com/vallard/drone-kube). The plugin will create a plugin if it doesn't already exist.

This plugin supports deployments, configmaps, ingresses, and services.

Usage: See [drone.md](./drone.md)

## Deployment, service, ingress templates

Deployment config files are first interpreted by **aymerick/raymond** ([handlebarsjs](http://handlebarsjs.com/) equivalent). You can use all available raymond expressions and anything you put in settings will be made available in your deployment template e.g. `{{namespace}}`. See [example/deployment.template.yaml](/example/deployment.template.yaml) for a complete example.

## Config maps from files

In this case, you can create a template just like deployment.yaml but you can provide a file path (relative to the repo's root) in the plugin setting `configmap_file`. (Currently only accepts utf-8 encoded data). Like deployments, this will both create new or update existing configmaps (based on the configmap name).

#### Adding a service account to Kubernetes that can manage deployments
See [example/Role.yaml](example/Role.yaml), [example/ServiceAccount.yaml](example/ServiceAccount.yaml), [example/RoleBinding.yaml](example/RoleBinding.yaml).

## Notes:

- The watching process is not currently reliable i.e. it doesn't properly wait for new deployments to become live. Not entirely sure how this should behave but I think behind a flag would make sense. PRs welcome.

## Development notes
- Kubernetes client is a little confusing with dependencies but does work with go.mod as seen [here](https://github.com/kubernetes/client-go/blob/master/INSTALL.md#add-client-go-as-a-dependency)
- [kubernetes/client-go installation notes](https://github.com/kubernetes/client-go/blob/master/INSTALL.md)
- [Creating a Drone plugin in Go](https://docs.drone.io/plugins/tutorials/golang/)
- [Client-go API Docs @ godoc.org](https://godoc.org/k8s.io/client-go/kubernetes)
- Testing with minikube (OSX: `brew cask install minikube`)

## Acknowledgements
- [@vallard](https://github.com/vallard) for the original plugin
- [@jbonzo](https://github.com/jbonzo) for ingress, service support, improved error handling etc
