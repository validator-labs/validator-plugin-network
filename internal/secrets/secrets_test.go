package secrets

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type mockReader struct {
	secretData map[string][]byte
}

func (m mockReader) Get(_ context.Context, key client.ObjectKey, obj client.Object, _ ...client.GetOption) error {
	switch o := obj.(type) {
	case *corev1.Secret:
		if m.secretData == nil {
			return fmt.Errorf("secret %s/%s not found", key.Namespace, key.Name)
		}
		o.Data = m.secretData
	}
	return nil
}

func (m mockReader) List(_ context.Context, _ client.ObjectList, _ ...client.ListOption) error {
	return nil
}

func TestReadKeys(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		keys      []string
		reader    client.Reader
		want      [][]byte
		wantErr   bool
	}{
		{
			name: "Pass",
			keys: []string{"key1", "key2"},
			reader: mockReader{
				secretData: map[string][]byte{
					"key1": []byte("value1"),
					"key2": []byte("value2"),
				},
			},
			want: [][]byte{[]byte("value1"), []byte("value2")},
		},
		{
			name: "Fail: secret does not contain key",
			keys: []string{"key1", "key2"},
			reader: mockReader{
				secretData: map[string][]byte{
					"key1": []byte("value1"),
				},
			},
			wantErr: true,
		},
		{
			name:    "Fail: secret not found",
			keys:    nil,
			reader:  mockReader{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadKeys(tt.name, tt.namespace, tt.keys, tt.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
