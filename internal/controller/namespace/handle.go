package namespace

import (
	"context"
	"fmt"

	"github.com/spotahome/kooper/v2/controller"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/slok/imagepull-controller-workshop/internal/log"
)

// HandlerRepository is the service to manage k8s resources by the Kubernetes controller handler.
type HandlerRepository interface {
	GetSecret(ctx context.Context, ns string, name string) (*corev1.Secret, error)
	EnsureSecret(ctx context.Context, secret *corev1.Secret) error
	GetServiceAccount(ctx context.Context, ns string, name string) (*corev1.ServiceAccount, error)
	EnsureServiceAccount(ctx context.Context, sa *corev1.ServiceAccount) error
}

// HandlerConfig is the handler configuration.
type HandlerConfig struct {
	RunningNamespace      string
	ImagePullSecretName   string
	SaImagePullSecretName string
	K8sRepo               HandlerRepository
	Logger                log.Logger
}

func (c *HandlerConfig) defaults() error {
	if c.RunningNamespace == "" {
		return fmt.Errorf("running namespaces is required")
	}

	if c.ImagePullSecretName == "" {
		c.ImagePullSecretName = "image-pull-credentials"
	}

	if c.SaImagePullSecretName == "" {
		c.SaImagePullSecretName = c.ImagePullSecretName
	}

	if c.K8sRepo == nil {
		return fmt.Errorf("kubernetes repository is required")
	}

	if c.Logger == nil {
		c.Logger = log.Noop
	}
	c.Logger = c.Logger.WithValues(log.Kv{"svc": "controller.namespace.Handler"})

	return nil
}

type handler struct {
	runningNamespace      string
	imagePullSecretName   string
	saImagePullSecretName string
	k8sRepo               HandlerRepository
	logger                log.Logger
}

// NewHandler returns the handler for the controller.
func NewHandler(config HandlerConfig) (controller.Handler, error) {
	err := config.defaults()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return handler{
		runningNamespace:      config.RunningNamespace,
		imagePullSecretName:   config.ImagePullSecretName,
		saImagePullSecretName: config.SaImagePullSecretName,
		k8sRepo:               config.K8sRepo,
		logger:                config.Logger,
	}, nil
}

func (h handler) Handle(ctx context.Context, obj runtime.Object) error {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		h.logger.Warningf("controller received object that is not a namespace")
		return nil
	}
	logger := h.logger.WithValues(log.Kv{"k8s-name": ns.Name})

	// If is our same namespace, ignore handling.
	if ns.Name == h.runningNamespace {
		return nil
	}

	// Make a copy just in case of global mutation.
	ns = ns.DeepCopy()

	logger.Infof("Handling namespace")

	// Get secret from running namespace with docker registry credentials.
	secret, err := h.k8sRepo.GetSecret(ctx, h.runningNamespace, h.imagePullSecretName)
	if err != nil {
		return fmt.Errorf("could not retrieve docker registry credentials secret: %w", err)
	}

	// Ensure secret on expected namespace.
	annotations := map[string]string{}
	for k, v := range secret.Annotations {
		annotations[k] = v
	}
	annotations["app.kubernetes.io/managed-by"] = "imagepull-controller-workshop"

	newNsSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        h.saImagePullSecretName,
			Namespace:   ns.Name,
			Labels:      secret.Labels,
			Annotations: annotations,
		},
		Data: secret.Data,
		Type: secret.Type,
	}

	err = h.k8sRepo.EnsureSecret(ctx, newNsSecret)
	if err != nil {
		return fmt.Errorf("could not ensure docker registry credentials secret on namespace: %w", err)
	}

	// Patch `default` service account on expected namespace.
	sa, err := h.k8sRepo.GetServiceAccount(ctx, ns.Name, "default")
	if err != nil {
		return fmt.Errorf("could not retrieve default service account from namespace: %w", err)
	}
	if containsLocalObjectRef(sa.ImagePullSecrets, h.saImagePullSecretName) {
		// Already set, move along.
		h.logger.Debugf("'default' service account image pull secet already set")
		return nil
	}

	sa.ImagePullSecrets = append(sa.ImagePullSecrets, corev1.LocalObjectReference{Name: h.saImagePullSecretName})
	err = h.k8sRepo.EnsureServiceAccount(ctx, sa)
	if err != nil {
		return fmt.Errorf("could not ensure default service account: %w", err)
	}

	return nil
}

func containsLocalObjectRef(refs []corev1.LocalObjectReference, name string) bool {
	for _, ref := range refs {
		if ref.Name == name {
			return true
		}
	}
	return false
}
