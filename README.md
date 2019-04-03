# drone-plugin-kube

[![](https://images.microbadger.com/badges/version/danielgormly/drone-plugin-kube.svg)](https://microbadger.com/images/danielgormly/drone-plugin-kube "Get your own version badge on microbadger.com")

A simple Drone plugin for managing Kubernetes deployments. Follows from [vallard/drone-kube](https://github.com/vallard/drone-kube) but with dependency management, up-to-date client-go, docs updated to Drone 1.0.0 syntax, examples and a different structure.

## Usage

Add the following [build step](https://docs.drone.io/user-guide/pipeline/steps/) to your drone pipeline definition.

#### drone.yaml partial example
```yml
- name: Deploy app
  image: danielgormly/drone-plugin-kube
  settings:
    template: path/to/deployment.yaml # within repo
    ca: LS0tLS1... # BASE64 encoded string of the K8s CA cert
    server: https://10.0.0.20:6443 # K8s master node address
    token:
      from_secret: kubernetes_token # Service account token to a service account that can manage deployments
    namespace: custom # [Optional] Custom namespace. (Defaults to `default`)
    custom: string # [Optional] Available to be referenced in template rendering as PLUGIN_CUSTOM
    master_alias: production # [Optional] Custom setting example. Available as PLUGIN_MASTER_ALIAS
```

## deployment templates

Deployment config files are first interpreted by **aymerick/raymond** ([handlebarsjs](http://handlebarsjs.com/) equivalent). You can use all available raymond expressions and anything you put in settings prefixed with the PLUGIN_* environment variables e.g. `{{PLUGIN.NAMESPACE}}`. See [example/deployment.template.yaml](/example/deployment.template.yaml) for a complete example.

#### Adding a service account to Kubernetes that can manage deployments
See [example/Role.yaml](example/Role.yaml), [example/ServiceAccount.yaml](example/ServiceAccount.yaml), [example/RoleBinding.yaml](example/RoleBinding.yaml).

## Development notes
- No tagged releases or support for older go libraries yet, happy to take feedback in Github issues or PRs.
- Kubernetes client not yet supported by dep, so we are using
[`brew install glide`](https://github.com/Masterminds/glide).
- Update dependencies with brew `glide update --strip-vendor`
- [Creating a Drone plugin in Go](https://docs.drone.io/plugins/examples/golang/)
- Testing with minikube (OSX: `brew cask install minikube`)
