// Package secrets contains helpers for managing secrets.
package secrets

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
)

// Read reads a secret from the Kubernetes API.
func Read(secretName, namespace, dataKey string, reader client.Reader) ([]byte, error) {
	secret := &corev1.Secret{}
	err := reader.Get(context.TODO(), client.ObjectKey{
		Namespace: namespace,
		Name:      secretName,
	}, secret)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read secret: %w", err)
	}
	data, ok := secret.Data[dataKey]
	if !ok {
		return []byte{}, fmt.Errorf("secret does not contain key '%s'", dataKey)
	}
	return data, nil
}
