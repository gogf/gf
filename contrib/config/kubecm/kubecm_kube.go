// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package kubecm

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/gogf/gf/v2/os/gfile"
)

const (
	defaultKubernetesUserAgent  = `kubecm.Client`
	kubernetesNamespaceFilePath = `/var/run/secrets/kubernetes.io/serviceaccount/namespace`
)

// Namespace retrieves and returns the namespace of current pod.
// Note that this function should be called in kubernetes pod.
func Namespace() string {
	return gfile.GetContents(kubernetesNamespaceFilePath)
}

// NewDefaultKubeClient creates and returns a default kubernetes client.
// It is commonly used when the service is running inside kubernetes cluster.
func NewDefaultKubeClient(ctx context.Context) (*kubernetes.Clientset, error) {
	return NewKubeClientFromPath(ctx, "")
}

// NewKubeClientFromPath creates and returns a kubernetes REST client by given `kubeConfigFilePath`.
func NewKubeClientFromPath(ctx context.Context, kubeConfigFilePath string) (*kubernetes.Clientset, error) {
	restConfig, err := NewKubeConfigFromPath(ctx, kubeConfigFilePath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

// NewKubeClientFromConfig creates and returns client by given `rest.Config`.
func NewKubeClientFromConfig(ctx context.Context, config *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(config)
}

// NewDefaultKubeConfig creates and returns a default kubernetes config.
// It is commonly used when the service is running inside kubernetes cluster.
func NewDefaultKubeConfig(ctx context.Context) (*rest.Config, error) {
	return NewKubeConfigFromPath(ctx, "")
}

// NewKubeConfigFromPath creates and returns rest.Config object from given `kubeConfigFilePath`.
func NewKubeConfigFromPath(ctx context.Context, kubeConfigFilePath string) (*rest.Config, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFilePath)
	if err != nil {
		return nil, err
	}
	restConfig.UserAgent = defaultKubernetesUserAgent
	return restConfig, nil
}
