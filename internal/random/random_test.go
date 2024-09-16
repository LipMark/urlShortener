package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomAlias(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 2",
			size: 2,
		},
		{
			name: "size = 4",
			size: 4,
		},
		{
			name: "size = 8",
			size: 8,
		},
		{
			name: "size = 16",
			size: 16,
		},
		{
			name: "size = 32",
			size: 32,
		},
		{
			name: "size = 64",
			size: 64,
		},
		{
			name: "size = 128",
			size: 128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomAlias(tt.size)
			str2 := NewRandomAlias(tt.size)

			// check both for given length
			assert.Len(t, str1, tt.size)
			assert.Len(t, str2, tt.size)

			// check difference
			assert.NotEqual(t, str1, str2)
		})
	}
}
