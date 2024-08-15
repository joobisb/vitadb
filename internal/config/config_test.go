package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("Load from file", func(t *testing.T) {
		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, int64(104857600), cfg.MaxLogSize)
		assert.Equal(t, cfg.DoAsyncRepair, false)
		assert.Equal(t, "/tmp/vitadb", cfg.WALDir)
	})

	t.Run("Load with env vars", func(t *testing.T) {
		os.Setenv("VITADB_MAX_LOG_SIZE", "314572800")
		os.Setenv("VITADB_DO_ASYNC_REPAIR", "false")
		os.Setenv("VITADB_WAL_DIR", "/tmp/vitadb_test")
		defer os.Unsetenv("VITADB_MAX_LOG_SIZE")
		defer os.Unsetenv("VITADB_DO_ASYNC_REPAIR")
		defer os.Unsetenv("VITADB_WAL_DIR")

		viper.Reset()
		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, int64(314572800), cfg.MaxLogSize)
		assert.False(t, cfg.DoAsyncRepair)
		assert.Equal(t, "/tmp/vitadb_test", cfg.WALDir)
	})

	t.Run("Load defaults", func(t *testing.T) {
		viper.Reset()
		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, int64(104857600), cfg.MaxLogSize)
		assert.False(t, cfg.DoAsyncRepair)
		assert.Equal(t, "/tmp/vitadb", cfg.WALDir)
	})
}
