package pllconfig

import (
	"testing"

	paths "github.com/komarovn654/embedded_configurator/config/common_paths"
	"github.com/stretchr/testify/require"
)

type TestTarget struct {
}

func (t *TestTarget) SetupPll() error {
	return nil
}

func TestNew(t *testing.T) {
	t.Run("new pllconfig", func(t *testing.T) {
		pc := New()
		require.NotNil(t, pc)
		require.NotNil(t, pc.Paths)
		require.Nil(t, pc.Target)
	})
}

func TestSetTarget(t *testing.T) {
	t.Run("set target", func(t *testing.T) {
		pc := PllConfig{}
		pc.SetTarget(&TestTarget{})
		require.Equal(t, &TestTarget{}, pc.Target)
	})
}

func TestGetTarget(t *testing.T) {
	tests := []struct {
		name     string
		target   PllTargetIf
		expected PllTargetIf
	}{
		{
			name:     "not nil target",
			target:   &TestTarget{},
			expected: &TestTarget{},
		},
		{
			name:     "nil target",
			target:   nil,
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run("get target", func(t *testing.T) {
			pc := PllConfig{Target: tc.target}
			require.Equal(t, tc.expected, pc.GetTarget())
		})
	}
}

func TestSetPath(t *testing.T) {
	t.Run("set path", func(t *testing.T) {
		p := paths.New()
		pc := PllConfig{}
		pc.SetPaths(p)
		require.Equal(t, p, pc.Paths)
	})
}

func TestGetPath(t *testing.T) {
	t.Run("get path", func(t *testing.T) {
		p := paths.New()
		pc := PllConfig{Paths: p}
		require.Equal(t, p, pc.GetPaths())
	})
}
