// Package secrets contains helpers for managing secrets.
package secrets

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
)

// ReadKeys fetches a Kubernetes secret and returns the value(s) of the specified key(s).
func ReadKeys(name, namespace string, keys []string, reader client.Reader) ([][]byte, error) {
	secret := &corev1.Secret{}
	if err := reader.Get(context.TODO(), client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, secret); err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}
	values := make([][]byte, len(keys))
	for i, k := range keys {
		data, ok := secret.Data[k]
		if !ok {
			return nil, fmt.Errorf("secret %s/%s does not contain key %s", name, namespace, k)
		}
		values[i] = data
	}
	return values, nil
}
