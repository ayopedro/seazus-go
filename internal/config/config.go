package config

import (
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv         string   `mapstructure:"APP_ENV"`
	Port           string   `mapstructure:"PORT"`
	LogLevel       string   `mapstructure:"LOG_LEVEL"`
	TrustedOrigins []string `mapstructure:"TRUSTED_ORIGINS"`
	EncryptionKey  string   `mapstructure:"ENCRYPTION_KEY"`
	JWTSecret      string   `mapstructure:"JWT_SECRET"`
	DB             DBConfig `mapstructure:",squash"`
	Limiter        RateLimitConfig
}

type DBConfig struct {
	User         string `mapstructure:"DB_USER"`
	Password     string `mapstructure:"DB_PASSWORD"`
	Host         string `mapstructure:"DB_HOST"`
	Port         string `mapstructure:"DB_PORT"`
	Name         string `mapstructure:"DB_NAME"`
	MaxOpenConns int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	MaxIdleTime  string `mapstructure:"DB_MAX_IDLE_TIME"`
}

type RateLimitConfig struct {
	RequestsPerTimeframe int           `mapstructure:"RATE_LIMIT_REQUEST_PER_TIMEFRAME"`
	Timeframe            time.Duration `mapstructure:"RATE_LIMIT_TIMEFRAME"`
	Enabled              bool          `mapstructure:"RATE_LIMIT_ENABLED"`
}

var (
	once   sync.Once
	config *Config
)

func Load() *Config {
	once.Do(func() {
		viper.SetDefault("APP_ENV", "development")
		viper.SetDefault("PORT", "8080")
		viper.SetDefault("DB_HOST", "localhost")
		viper.SetDefault("DB_PORT", "5432")
		viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
		viper.SetDefault("DB_MAX_IDLE_CONNS", 25)
		viper.SetDefault("DB_MAX_IDLE_TIME", "15m")
		viper.SetDefault("LOG_LEVEL", "info")
		viper.SetDefault("RATE_LIMIT_REQUEST_PER_TIMEFRAME", 100)
		viper.SetDefault("RATE_LIMIT_TIMEFRAME", 1*time.Minute)
		viper.SetDefault("TRUSTED_ORIGINS", []string{"*"})
		viper.SetDefault("JWT_SECRET", "my-super-not-so-secret-secret")
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
