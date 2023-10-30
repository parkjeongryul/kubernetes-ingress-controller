package manager_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/manager"
)

func TestConfigResolve(t *testing.T) {
	t.Run("Admin Token Path", func(t *testing.T) {
		validWithTokenPath := func() manager.Config {
			tempDir := t.TempDir()
			tokenFile, err := os.CreateTemp(tempDir, "kong.token")
			require.NoError(t, err)
			_, err = tokenFile.Write([]byte("non-empty-token"))
			require.NoError(t, err)
			return manager.Config{
				KongAdminTokenPath: tokenFile.Name(),
			}
		}

		t.Run("admin token path accepted", func(t *testing.T) {
			c := validWithTokenPath()
			require.NoError(t, c.Resolve())
			require.Equal(t, c.KongAdminToken, "non-empty-token")
		})
	})
}
