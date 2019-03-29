# kubano

A simple Drone plugin for managing Kubernetes deployments.

## Usage

Add the following [build step](https://docs.drone.io/user-guide/pipeline/steps/) to your drone pipeline definition.

#### drone.yaml partial example
```yml
- name:
  image: danielgormly/kubano
  settings:
    template: deployment.yaml
```

#### Settings
- **template**: path to deployment file from repo root e.g. `deploy/deployment.yaml`, etc

## deployment templates

Deployment config files are first interpreted by raymond ([handlebarsjs](http://handlebarsjs.com/) equivalent). Use `{{variable}}` to add interpolated expressions e.g.

#### deployment.yaml partial example
```yaml
spec:
  containers:
  - name: nginx
    image: 10.0.0.24:443/danielgormly:{{git-repo}}.{{git-branch}}
```

#### Development
- Kubernetes client not yet supported by dep, so we are using
[`brew install glide`](https://github.com/Masterminds/glide).
- Update dependencies with brew `glide update --strip-vendor`
- [Creating a Drone plugin in Go](https://docs.drone.io/plugins/examples/golang/)

#### Acknowledgements
- Heavily referenced [vallard/drone-kube](https://github.com/vallard/drone-kube).
