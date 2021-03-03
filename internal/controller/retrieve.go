package controller

import (
	"context"

	"github.com/spotahome/kooper/v2/controller"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/slok/imagepull-controller-workshop/internal/log"
)

// RetrieverKubernetesRepository is the service to manage k8s resources by the Kubernetes retrievers.
type RetrieverKubernetesRepository interface {
	ListNamespaces(ctx context.Context, labelSelector map[string]string) (*corev1.NamespaceList, error)
	WatchNamespaces(ctx context.Context, labelSelector map[string]string) (watch.Interface, error)
}

// NewRetriever returns the retriever for the controller.
func NewRetriever(k8sRepo RetrieverKubernetesRepository, logger log.Logger) (controller.Retriever, error) {
	return controller.RetrieverFromListerWatcher(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return k8sRepo.ListNamespaces(context.Background(), map[string]string{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return k8sRepo.WatchNamespaces(context.Background(), map[string]string{})
		},
	})
}
