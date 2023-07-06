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

func TestSetTemplatePath(t *testing.T) {
	t.Run("set path", func(t *testing.T) {
		p := Paths{}
		option := SetTemplatePath("new path")
		require.IsType(t, func(*Paths) {}, option)

		option(&p)
		require.Equal(t, "new path", p.TemplatePath)
		require.Equal(t, "", p.DestinationPath)
	})

	t.Run("set path via class method", func(t *testing.T) {
		p := Paths{}
		p.SetTemplatePath("new path")
		require.Equal(t, "new path", p.TemplatePath)
		require.Equal(t, "", p.DestinationPath)
	})
}

func TestSetDestinationPath(t *testing.T) {
	t.Run("set path", func(t *testing.T) {
		p := Paths{}
		option := SetDestinationPath("new path")
		require.IsType(t, func(*Paths) {}, option)

		option(&p)
		require.Equal(t, "new path", p.DestinationPath)
		require.Equal(t, "", p.TemplatePath)
	})

	t.Run("set path via class method", func(t *testing.T) {
		p := Paths{}
		p.SetDestinationPath("new path")
		require.Equal(t, "new path", p.DestinationPath)
		require.Equal(t, "", p.TemplatePath)
	})
}

func TestGetTemplatePath(t *testing.T) {
	t.Run("get path", func(t *testing.T) {
		p := Paths{TemplatePath: "new path"}
		require.Equal(t, "new path", p.GetTemplatePath())
	})
}

func TestGetDestinationPath(t *testing.T) {
	t.Run("get path", func(t *testing.T) {
		p := Paths{DestinationPath: "new path"}
		require.Equal(t, "new path", p.GetDestinationPath())
	})
}
