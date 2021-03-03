package kubernetes

import (
	"context"
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
)

// SecretCachedRepository is a Kubernetes repository like `Repository` but getting a
// secret is done from an internal cache.
type SecretCachedRepository struct {
	mu      sync.RWMutex
	secrets map[string]*corev1.Secret
	Repository
}

// NewSecretCachedRepository returns a new NewSecretCachedRepository.
func NewSecretCachedRepository(repo Repository) *SecretCachedRepository {
	return &SecretCachedRepository{
		Repository: repo,
		secrets:    map[string]*corev1.Secret{},
	}
}

// GetSecret will return a secret from the internal cache.
func (s *SecretCachedRepository) GetSecret(ctx context.Context, ns string, name string) (*corev1.Secret, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id := fmt.Sprintf("%s/%s", ns, name)
	secret, ok := s.secrets[id]
	if !ok {
		return nil, fmt.Errorf("secret not found")
	}

	// Deep copy because we don't want cache object mutations from the outside.
	return secret.DeepCopy(), nil
}

// SetSecretOnCache sets a new secret on the repository cache.
func (s *SecretCachedRepository) SetSecretOnCache(ctx context.Context, secret *corev1.Secret) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("%s/%s", secret.Namespace, secret.Name)
	s.secrets[id] = secret

	return nil
}
