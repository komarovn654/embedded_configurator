package paths

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		options      []Option
		expectedTmpl string
		expectedDst  string
	}{
		{
			name:         "no options",
			options:      nil,
			expectedTmpl: "",
			expectedDst:  "",
		},
		{
			name:         "with options",
			options:      []Option{SetTemplatePath("tmplt"), SetDestinationPath("dst")},
			expectedTmpl: "tmplt",
			expectedDst:  "dst",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := New(tc.options...)
			require.Equal(t, tc.expectedTmpl, p.TemplatePath)
			require.Equal(t, tc.expectedDst, p.DestinationPath)
		})
	}
}
