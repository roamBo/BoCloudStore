package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Env        string
	ServerPort string
	Minio      MinioConfig
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	viper.SetDefault("env", "development")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("minio.endpoint", "minio:9000")
	viper.SetDefault("minio.accessKey", "minioadmin")
	viper.SetDefault("minio.secretKey", "minioadmin")
	viper.SetDefault("minio.useSSL", false)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}

	return &Config{
		Env:        viper.GetString("env"),
		ServerPort: viper.GetString("server.port"),
		Minio: MinioConfig{
			Endpoint:  viper.GetString("minio.endpoint"),
			AccessKey: viper.GetString("minio.accessKey"),
			SecretKey: viper.GetString("minio.secretKey"),
			UseSSL:    viper.GetBool("minio.useSSL"),
		},
	}
}
