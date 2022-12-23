package config

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	App struct {
		Name       string
		Version    string
		Host       string
		Port       string
		CtxTimeout time.Duration
	}
	Database Database
}

type Database struct {
	DebugLog            bool
	SSL                 string
	Host                string
	Port                string
	User                string
	Pass                string
	DatabaseName        string
	Charset             string
	Timezone            string
	PoolMaxIdleConn     int
	PoolMaxOpenConn     int
	PoolMaxConnLifetime time.Duration
}

func LoadConfigFile(path string) *Config {
	viperConfig := viper.New()
	viperConfig.SetConfigName(path)
	viperConfig.SetConfigType("yaml")
	viperConfig.AddConfigPath(".")

	if err := viperConfig.ReadInConfig(); err != nil {
		log.Fatalf("failed reading config, %v", err)
	}

	config := new(Config)
	if err := viperConfig.Unmarshal(&config); err != nil {
		log.Fatalf("failed parsing config, %v", err)
	}

	log.Printf("config loaded")

	return config
}
