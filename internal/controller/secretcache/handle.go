package secretcache

import (
	"context"
	"fmt"

	"github.com/spotahome/kooper/v2/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/slok/imagepull-controller-workshop/internal/log"
)

// HandlerRepository is the service to manage k8s resources by the Kubernetes controller handler.
type HandlerRepository interface {
	SetSecretOnCache(ctx context.Context, secret *corev1.Secret) error
}

type handler struct {
	k8sRepo HandlerRepository
	logger  log.Logger
}

// NewHandler returns the handler for the controller.
func NewHandler(k8sRepo HandlerRepository, logger log.Logger) (controller.Handler, error) {
	return handler{
		k8sRepo: k8sRepo,
		logger:  logger.WithValues(log.Kv{"svc": "controller.secretcache.Handler"}),
	}, nil
}

func (h handler) Handle(ctx context.Context, obj runtime.Object) error {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		h.logger.Warningf("controller received object that is not a secret")
		return nil
	}

	logger := h.logger.WithValues(log.Kv{"k8s-ns": secret.Namespace, "k8s-name": secret.Name})

	// Store on cache.
	err := h.k8sRepo.SetSecretOnCache(ctx, secret)
	if err != nil {
		return fmt.Errorf("could not update secret cache: %w", err)
	}

	logger.Infof("Secret cache updated")

	return nil
}
