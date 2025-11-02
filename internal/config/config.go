package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		APIPrefix string `mapstructure:"api_prefix"`
		Port      int    `mapstructure:"port"`
	} `mapstructure:"server"`

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
