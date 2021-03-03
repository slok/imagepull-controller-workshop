package secretcache

import (
	"context"
	"fmt"

	"github.com/spotahome/kooper/v2/controller"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// RetrieverRepository is the service to manage k8s resources by the Kubernetes retrievers.
type RetrieverRepository interface {
	ListSecrets(ctx context.Context, ns string, options metav1.ListOptions) (*corev1.SecretList, error)
	WatchSecrets(ctx context.Context, ns string, options metav1.ListOptions) (watch.Interface, error)
}

// NewRetriever returns the retriever for the controller.
func NewRetriever(k8sRepo RetrieverRepository, ns, secretName string) (controller.Retriever, error) {
	secretFieldSelector := fmt.Sprintf("metadata.name=%s", secretName)

	return controller.RetrieverFromListerWatcher(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = secretFieldSelector
			return k8sRepo.ListSecrets(context.Background(), ns, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = secretFieldSelector
			return k8sRepo.WatchSecrets(context.Background(), ns, options)
		},
	})
}
