### DynRateLimiter Operator

Create your `values.yaml`

```yaml
# your values.yaml to override default values
namespace: finnie-pipeline-operator-system
image:
	pullSecrets:
		- name: my-pullsecret
```

Then, install the `DynRateLimiter Operator` with Helm:

```bash
# deploy using your values.yaml
helm install -f values.yaml dynratelimiter-operator arivum/dynratelimiter-operator
```

#### Available chart properties

| Property       |   Description         | Default               |
| -------------- | --------------------- | --------------------- |
| `replicaCount` | Number of replicas    | `1`                   |
| `appname`      | App Name, certificates and service name will be derived from this value    | `dynratelimiter-operator`      |
| `namespace`    | Namespace to deploy the operator to  | `dynratelimiter-operator`  |
| `image.repository`  | Image repository | `ghcr.io/arivum/dynratelimiter/dynratelimiter-operator` |
| `image.tag`  | Image tag  | `latest` |
| `image.pullPolicy`  | Image pull policy | `IfNotPresent` |
| `imagePullSecrets`  | A list of image pull secrets when using your own image | `[]` |
| `nameOverride`  | Overrides the basename for deployments etc. | `""` |
| `fullnameOverride` | Completly overrides the auto-generated names for deployments etc. | `""` |
| `logging.level` | Log level. Must be one of `[info, debug, warn, error, trace]` | `info` |
| `logging.format` | Log format. Choose one of `[gofmt, json]` | `gofmt` |
| `mutationConfig.image` | Configures the image that gets injected into annotated pods by the operator | `ghcr.io/arivum/dynratelimiter/dynratelimiter` |
| `mutationConfig.tag` | Configures the tag of image that gets injected into annotated pods by the operator | `latest` |
| `service.type` | Kubernetes serivce type | `ClusterIP` |
| `service.port` | Kubernetes service port | `443` |
| `service.targetPort` | Kubernetes service target port and listen port of the dynratelimit-operator application | `8443` |
| `resources` | Pod resource requests and limits | `{}` |
| `nodeSelector` | Kubernetes node selector | `{}` |
| `tolerations` | Kuberentes tolerations | `[]` |
| `affinity` | Kubernetes pod affinity | `{}` |