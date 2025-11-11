package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		APIPrefix string `mapstructure:"api_prefix"`
		Port      int    `mapstructure:"port"`
	} `mapstructure:"server"`

	JWT struct {
		AccessName       string        `mapstructure:"access_name"`
		RefreshName      string        `mapstructure:"refresh_name"`
		SecretKey        string        `mapstructure:"secret_key"`
		AccessExpiresIn  time.Duration `mapstructure:"access_expires_in"`
		RefreshExpiresIn time.Duration `mapstructure:"refresh_expires_in"`
	}

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`

	RabbitMQ struct {
		Host     string `mapstructure:"host"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Vhost    string `mapstructure:"vhost"`
	} `mapstructure:"rabbitmq"`

	PostgreSQL struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		SSLMode  string `mapstructure:"ssl_mode"`
		DBName   string `mapstructure:"db_name"`
	} `mapstructure:"postgresql"`

	S3 struct {
		Bucket          string `mapstructure:"bucket"`
		Folder          string `mapstructure:"folder"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		Region          string `mapstructure:"region"`
	} `mapstructure:"s3"`

	SMTP struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"smtp"`

	IMAP struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"imap"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("configs/main.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
