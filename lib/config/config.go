package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	Port           string `mapstructure:"PORT"`
	DBURL          string `mapstructure:"DATABASE_URL"`
	DBMaxOpenConns int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxIdleTime  string `mapstructure:"DB_MAX_IDLE_TIME"`
	LogLevel       string `mapstructure:"LOG_LEVEL"`
}

var (
	once   sync.Once
	config *Config
)

func Load() *Config {
	once.Do(func() {
		viper.SetDefault("APP_ENV", "development")
		viper.SetDefault("PORT", "8080")
		viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
		viper.SetDefault("DB_MAX_IDLE_CONNS", 25)
		viper.SetDefault("DB_MAX_IDLE_TIME", "15m")
		viper.SetDefault("LOG_LEVEL", "info")

		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatalf("Error reading config file: %v", err)
			}
		}

		config = &Config{}
		if err := viper.Unmarshal(config); err != nil {
			log.Fatalf("Unable to decode into struct: %v", err)
		}
	})

	return config
}
