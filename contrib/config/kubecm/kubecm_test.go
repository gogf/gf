// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package kubecm_test

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	kubeMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/gogf/gf/contrib/config/kubecm/v2"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

const (
	namespace         = "default"
	configmap         = "test-configmap"
	dataItem          = "config.yaml"
	configmapFileName = "configmap.yaml"
)

var (
	ctx                    = gctx.New()
	kubeConfigFilePath     = `/home/runner/.kube/config`
	kubeConfigFilePathJohn = `/Users/john/.kube/config`
)

func init() {
	if !gfile.Exists(kubeConfigFilePath) {
		kubeConfigFilePath = kubeConfigFilePathJohn
	}
}

func TestAvailable(t *testing.T) {
	var (
		err        error
		kubeClient *kubernetes.Clientset
	)
	// Configmap apply.
	gtest.C(t, func(t *gtest.T) {
		kubeClient, err = kubecm.NewKubeClientFromPath(ctx, kubeConfigFilePath)
		t.AssertNil(err)
		var (
			configMap v1.ConfigMap
			content   = gtest.DataContent(configmapFileName)
		)
		err = gjson.New(content).Scan(&configMap)
		t.AssertNil(err)
		_, err = kubeClient.CoreV1().ConfigMaps(namespace).Create(
			ctx, &configMap, kubeMetaV1.CreateOptions{},
		)
		t.AssertNil(err)
	})
	defer func() {
		gtest.C(t, func(t *gtest.T) {
			err = kubeClient.CoreV1().ConfigMaps(namespace).Delete(
				ctx, configmap, kubeMetaV1.DeleteOptions{},
			)
			t.AssertNil(err)
		})
	}()

	gtest.C(t, func(t *gtest.T) {
		adapter, err := kubecm.New(ctx, kubecm.Config{
			ConfigMap:  configmap,
			DataItem:   dataItem,
			Namespace:  namespace,
			KubeClient: kubeClient,
		})
		t.AssertNil(err)

		config := g.Cfg(guid.S())
		config.SetAdapter(adapter)

		t.Assert(config.Available(ctx), true)

		m, err := config.Data(ctx)
		t.AssertNil(err)
		t.AssertGT(len(m), 0)

		v, err := config.Get(ctx, "server.address")
		t.AssertNil(err)
		t.Assert(v.String(), ":8888")
	})
}

func TestNewKubeClientFromConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config, _ := kubecm.NewKubeConfigFromPath(ctx, kubeConfigFilePath)
		_, err := kubecm.NewKubeClientFromConfig(ctx, config)
		t.AssertNil(err)
	})
}

// These functions should be called in pod environment, but it has no environment in CI UT testing.
// It so just calls them ,but does nothing.
func TestDefaultBehaviorFunctions(t *testing.T) {
	kubecm.Namespace()
	kubecm.NewDefaultKubeClient(ctx)
	kubecm.NewDefaultKubeConfig(ctx)
}
