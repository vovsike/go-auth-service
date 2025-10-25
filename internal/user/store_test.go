package user

import (
	"testing"
)

func TestNewInMemStore(t *testing.T) {
	tests := []struct {
		name string
		want InMemStore
	}{
		{
			name: "new in memory store",
			want: InMemStore{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInMemStore(); got == nil {
				t.Errorf("NewInMemStore() was not able to create a new in memory store")
			}
		})
	}
}
