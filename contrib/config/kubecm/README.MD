# kubecm
Package `kubecm` implements GoFrame `gcfg.Adapter` using kubernetes configmap.

# Limit

```go
glang version >= v.18
```

# Installation
```
go get -u github.com/gogf/gf/contrib/config/kubecm/v2
```

# Example

## Example configmap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-configmap
data:
  config.yaml: |
    # HTTP service.
    server:
      address:          ":8888"
      openapiPath:      "/api.json"
      swaggerPath:      "/swagger"
      accessLogEnabled: true
```

Note the configmap name `test-configmap`, and its item name in data `config.yaml`.


## Create a custom boot package

It is strongly recommended creating a boot package, 
which sets the Adapter of default configuration instance.

### Running in pod (common scenario)
```go
package boot

import (
	"github.com/gogf/gf/contrib/config/kubecm/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	configmapName       = "test-configmap"
	dataItemInConfigmap = "config.yaml"
)

func init() {
	var (
		err error
		ctx = gctx.GetInitCtx()
	)
	// Create kubecm Client that implements gcfg.Adapter.
	adapter, err := kubecm.New(gctx.GetInitCtx(), kubecm.Config{
		ConfigMap: configmapName,
		DataItem:  dataItemInConfigmap,
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}

	// Change the adapter of default configuration instance.
	g.Cfg().SetAdapter(adapter)
}
```

### Running out of pod (often testing scenario)
```go
package boot

import (
	"github.com/gogf/gf/contrib/config/kubecm/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"k8s.io/client-go/kubernetes"
)

const (
	namespace              = "default"
	configmapName          = "test-configmap"
	dataItemInConfigmap    = "config.yaml"
	kubeConfigFilePathJohn = `/Users/john/.kube/config`
)

func init() {
	var (
		err        error
		ctx        = gctx.GetInitCtx()
		kubeClient *kubernetes.Clientset
	)
	// Create kubernetes client.
	// It is optional creating kube client when its is running in pod.
	kubeClient, err = kubecm.NewKubeClientFromPath(ctx, kubeConfigFilePathJohn)
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}
	// Create kubecm Client that implements gcfg.Adapter.
	adapter, err := kubecm.New(gctx.GetInitCtx(), kubecm.Config{
		ConfigMap:  configmapName,
		DataItem:   dataItemInConfigmap,
		Namespace:  namespace,  // It is optional specifying namespace when its is running in pod.
		KubeClient: kubeClient, // It is optional specifying kube client when its is running in pod.
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}

	// Change the adapter of default configuration instance.
	g.Cfg().SetAdapter(adapter)

}
```

## Import boot package in top of main

It is strongly recommended import your boot package in top of your `main.go`.

Note the top `import`: `_ "github.com/gogf/gf/example/config/kubecm/boot_in_pod"` .

```go
package main

import (
	_ "github.com/gogf/gf/example/config/kubecm/boot_in_pod"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var ctx = gctx.GetInitCtx()

	// Available checks.
	g.Dump(g.Cfg().Available(ctx))

	// All key-value configurations.
	g.Dump(g.Cfg().Data(ctx))

	// Retrieve certain value by key.
	g.Dump(g.Cfg().MustGet(ctx, "server.address"))
}

```

