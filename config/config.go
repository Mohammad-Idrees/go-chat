package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type StartupConfig struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DBConfig        `mapstructure:"database"`
	Migration MigrationConfig `mapstructure:"migration"`
	Token     TokenConfig     `mapstructure:"token"`
	Redis     RedisConfig     `mapstructure:"redis"`
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
}

type DBConfig struct {
	Type     string `mapstructure:"type"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
}

type MigrationConfig struct {
	MigrationURL string `mapstructure:"migrationURL"`
}

type TokenConfig struct {
	JWTSecret            string        `mapstructure:"jwtSecret"`
	AccessTokenDuration  time.Duration `mapstructure:"accessTokenDuration"`
	RefreshTokenDuration time.Duration `mapstructure:"refreshTokenDuration"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

func LoadConfig() (*StartupConfig, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("failed reading config", err)
		return nil, err
	}

	var config StartupConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
