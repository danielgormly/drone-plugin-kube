# kubano

A simple Drone plugin for managing Kubernetes deployments.

## Usage

Add the following [build step](https://docs.drone.io/user-guide/pipeline/steps/) to your drone pipeline definition.

#### drone.yaml partial example
```yml
- name: Deploy app
  image: danielgormly/kubano
  settings:
    template: path/to/deployment.yaml # within repo
    ca: # BASE64 encoded string of the K8s CA cert
    Endpoint: 10.0.0.24 # K8s master node address
    Token: # Service account token to a service account that can manage deployments
    Namespace: custom # Custom namespace. (Optional, defaults to `default`)
    custom: string # Available to be referenced in template rendering as PLUGIN_CUSTOM
    master_alias: production # Custom setting example. Available as PLUGIN_MASTER_ALIAS
```

## deployment templates

Deployment config files are first interpreted by raymond ([handlebarsjs](http://handlebarsjs.com/) equivalent). You can use all available Use `{{VARIABLE}}` to add interpolated expressions e.g.

#### deployment.yaml partial example
```yaml
spec:
  containers:
  - name: nginx
    image: 10.0.0.24:443/danielgormly:{{DRONE_BRANCH}}
```

#### Development
- Kubernetes client not yet supported by dep, so we are using
[`brew install glide`](https://github.com/Masterminds/glide).
- Update dependencies with brew `glide update --strip-vendor`
- [Creating a Drone plugin in Go](https://docs.drone.io/plugins/examples/golang/)

#### Acknowledgements
- Heavily referenced [vallard/drone-kube](https://github.com/vallard/drone-kube).
