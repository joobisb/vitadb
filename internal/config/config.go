package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MaxLogSize    int64  `mapstructure:"max_log_size"`
	DoAsyncRepair bool   `mapstructure:"do_async_repair"`
	WALDir        string `mapstructure:"wal_dir"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	//relative path to config.yaml from cmd/main.go
	viper.AddConfigPath("../")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("VITADB")

	viper.SetDefault("max_log_size", 104857600)
	viper.SetDefault("do_async_repair", false)
	viper.SetDefault("wal_dir", "/tmp/vitadb")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			return nil, err
		}
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
