package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredential(t *testing.T) {
	tests := []struct {
		name string
		cred string
		want map[string]string
	}{
		{
			name: "user with pass",
			cred: "user:pass",
			want: map[string]string{"user": "pass"},
		}, {
			name: "user with pass 2",
			cred: "user:pass:123",
			want: map[string]string{"user": "pass:123"},
		}, {
			name: "pass only",
			cred: "pass",
			want: map[string]string{"fakecast": "pass"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := credential(tt.cred)
			assert.Equal(t, tt.want, got)
		})
	}
}
