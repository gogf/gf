package boot

import (
	"k8s.io/client-go/kubernetes"

	"github.com/gogf/gf/contrib/config/kubecm/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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
