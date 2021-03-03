package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// Repository represents a Kubernetes repository that knows how to speak with the
// Kubernetes API server to manage resources.
type Repository struct {
	kcli *kubernetes.Clientset
}

// NewRepository returns a new Kubernetes repository that will retrieve Kubernetes resources
// using kubernetes sdk.
func NewRepository(kcli *kubernetes.Clientset) Repository {
	return Repository{kcli: kcli}
}

// ListNamespaces will list Kubernetes namespaces from the API server.
func (r Repository) ListNamespaces(ctx context.Context, labelSelector map[string]string) (*corev1.NamespaceList, error) {
	return r.kcli.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector).String(),
	})
}

// WatchNamespaces will return a Kubernetes watcher to subscribe to namespaces changes.
func (r Repository) WatchNamespaces(ctx context.Context, labelSelector map[string]string) (watch.Interface, error) {
	return r.kcli.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector).String(),
	})
}

// GetSecret will return a secret from Kubernets API server.
func (r Repository) GetSecret(ctx context.Context, ns string, name string) (*corev1.Secret, error) {
	return r.kcli.CoreV1().Secrets(ns).Get(ctx, name, metav1.GetOptions{})
}

// ListSecrets lists Kubernetes secrets from Kubernetes API server.
func (r Repository) ListSecrets(ctx context.Context, ns string, options metav1.ListOptions) (*corev1.SecretList, error) {
	return r.kcli.CoreV1().Secrets(ns).List(ctx, options)
}

// WatchSecrets watchs Kubernetes secrets from Kubernetes API server.
func (r Repository) WatchSecrets(ctx context.Context, ns string, options metav1.ListOptions) (watch.Interface, error) {
	return r.kcli.CoreV1().Secrets(ns).Watch(ctx, options)
}

// EnsureSecret will create the secret if is missing and overwrite if already exists.
func (r Repository) EnsureSecret(ctx context.Context, secret *corev1.Secret) error {
	storedSecret, err := r.kcli.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{})
	if err != nil {
		if !kubeerrors.IsNotFound(err) {

			return err
		}

		_, err = r.kcli.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		return nil
	}

	// Force overwrite.
	secret.ObjectMeta.ResourceVersion = storedSecret.ResourceVersion
	_, err = r.kcli.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// GetServiceAccount  will return a service account from Kubernets API server.
func (r Repository) GetServiceAccount(ctx context.Context, ns string, name string) (*corev1.ServiceAccount, error) {
	return r.kcli.CoreV1().ServiceAccounts(ns).Get(ctx, name, metav1.GetOptions{})
}

// EnsureServiceAccount will create the service account if is missing and overwrite if already exists.
func (r Repository) EnsureServiceAccount(ctx context.Context, sa *corev1.ServiceAccount) error {
	storedSA, err := r.kcli.CoreV1().ServiceAccounts(sa.Namespace).Get(ctx, sa.Name, metav1.GetOptions{})
	if err != nil {
		if !kubeerrors.IsNotFound(err) {
			return err
		}

		_, err = r.kcli.CoreV1().ServiceAccounts(sa.Namespace).Create(ctx, sa, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		return nil
	}

	// Force overwrite.
	sa.ObjectMeta.ResourceVersion = storedSA.ResourceVersion
	_, err = r.kcli.CoreV1().ServiceAccounts(sa.Namespace).Update(ctx, sa, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
