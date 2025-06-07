package env_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xdrm-io/env"
)

func TestRead(t *testing.T) {
	f, err := os.CreateTemp("", "env_test")
	require.NoError(t, err)
	var (
		path    = f.Name()
		content = []byte("some content")
	)
	defer os.Remove(path)
	require.NoError(t, os.WriteFile(path, content, 0644))

	tt := []struct {
		name   string
		key    string
		envs   map[string]string
		expect string
		ok     bool
	}{
		{
			name:   "env ok",
			key:    "SIMPLE_KEY",
			envs:   map[string]string{"SIMPLE_KEY": "some value"},
			expect: "some value",
			ok:     true,
		},
		{
			name:   "env unset",
			key:    "SIMPLE_KEY",
			envs:   map[string]string{},
			expect: "",
			ok:     false,
		},
		{
			name:   "env empty",
			key:    "SIMPLE_KEY",
			envs:   map[string]string{"SIMPLE_KEY": ""},
			expect: "",
			ok:     true,
		},
		{
			name:   "file ok",
			key:    "SIMPLE_KEY",
			envs:   map[string]string{"SIMPLE_KEY_FILE": path},
			expect: string(content),
			ok:     true,
		},
		{
			name:   "file unset",
			key:    "SIMPLE_KEY",
			envs:   map[string]string{},
			expect: "",
			ok:     false,
		},
		{
			name:   "file empty",
			key:    "SIMPLE_KEY",
			envs:   map[string]string{"SIMPLE_KEY_FILE": ""},
			expect: "",
			ok:     false,
		},
		{
			name:   "file not found",
			key:    "SIMPLE_KEY",
			envs:   map[string]string{"SIMPLE_KEY_FILE": "/wrong/path"},
			expect: "",
			ok:     false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tc.envs {
				os.Setenv(k, v)
			}

			got, ok := env.Read(tc.key)
			require.Equal(t, tc.expect, got)
			require.Equal(t, tc.ok, ok)
		})
	}
}
