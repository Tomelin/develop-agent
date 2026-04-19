package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	App      AppConfig
	Mongo    MongoConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	LLM      LLMConfig
}

// AppConfig holds application specific configurations.
type AppConfig struct {
	Env   string
	Port  string
	Debug bool
}

// MongoConfig holds MongoDB configurations.
type MongoConfig struct {
	URI string
}

// RedisConfig holds Redis configurations.
type RedisConfig struct {
	Addr     string
	Password string
}

// RabbitMQConfig holds RabbitMQ configurations.
type RabbitMQConfig struct {
	URL string
}

// LLMConfig holds LLM configurations.
type LLMConfig struct {
	Provider string
	APIKey   string
}

// Load loads the configuration from a file and environment variables.
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Allow overriding configs with APP_ prefix via environment variables
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("app.debug", false)
	viper.SetDefault("mongo.uri", "mongodb://admin:password@localhost:27017")
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("rabbitmq.url", "amqp://guest:guest@localhost:5672/")

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// File found but contains errors
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// It's okay if the config file doesn't exist, we will use defaults + env vars
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err := validate(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func validate(cfg *Config) error {
	if cfg.App.Port == "" {
		return fmt.Errorf("missing critical configuration: app.port")
	}
	if cfg.Mongo.URI == "" {
		return fmt.Errorf("missing critical configuration: mongo.uri")
	}
	if cfg.Redis.Addr == "" {
		return fmt.Errorf("missing critical configuration: redis.addr")
	}
	if cfg.RabbitMQ.URL == "" {
		return fmt.Errorf("missing critical configuration: rabbitmq.url")
	}
	return nil
}
