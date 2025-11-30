package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		APIPrefix        string        `mapstructure:"api_prefix"`
		Port             int           `mapstructure:"port"`
		WriteTimeout     time.Duration `mapstructure:"write_timeout"`
		ReadTimeout      time.Duration `mapstructure:"read_timeout"`
		IdleTimeout      time.Duration `mapstructure:"idle_timeout"`
		MaxHeaderBytes   int           `mapstructure:"max_header_bytes"`
		AllowOrigins     []string      `mapstructure:"allow_origins"`
		AllowMethods     []string      `mapstructure:"allow_methods"`
		AllowHeaders     []string      `mapstructure:"allow_headers"`
		ExposeHeaders    []string      `mapstructure:"expose_headers"`
		AllowCredentials bool          `mapstructure:"allow_credentials"`
		MaxAge           time.Duration `mapstructure:"max_age"`
	} `mapstructure:"server"`

	JWT struct {
		AccessName       string        `mapstructure:"access_name"`
		RefreshName      string        `mapstructure:"refresh_name"`
		GuestName        string        `mapstructure:"guest_name"`
		SecretKey        string        `mapstructure:"secret_key"`
		AccessExpiresIn  time.Duration `mapstructure:"access_expires_in"`
		RefreshExpiresIn time.Duration `mapstructure:"refresh_expires_in"`
	} `mapstructure:"jwt"`

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
	viper.AutomaticEnv()
	viper.BindEnv("postgresql.host", "PG_HOST")
	viper.BindEnv("postgresql.port", "PG_PORT")
	viper.BindEnv("postgresql.user", "PG_USER")
	viper.BindEnv("postgresql.password", "PG_PASSWORD")
	viper.BindEnv("postgresql.ssl_mode", "PG_SSL_MODE")
	viper.BindEnv("postgresql.db_name", "PG_DB_NAME")

	viper.BindEnv("redis.host", "RD_HOST")
	viper.BindEnv("redis.port", "RD_PORT")
	viper.BindEnv("redis.password", "RD_PASSWORD")

	viper.BindEnv("rabbitmq.host", "RMQ_HOST")
	viper.BindEnv("rabbitmq.user", "RMQ_USER")
	viper.BindEnv("rabbitmq.password", "RMQ_PASSWORD")
	viper.BindEnv("rabbitmq.vhost", "RMQ_VHOST")

	viper.BindEnv("s3.bucket", "S3_BUCKET")
	viper.BindEnv("s3.folder", "S3_FOLDER")
	viper.BindEnv("s3.access_key_id", "S3_ACCESS_KEY_ID")
	viper.BindEnv("s3.secret_access_key", "S3_SECRET_ACCESS_KEY")
	viper.BindEnv("s3.region", "S3_REGION")

	viper.BindEnv("smtp.host", "SMTP_HOST")
	viper.BindEnv("smtp.port", "SMTP_PORT")
	viper.BindEnv("smtp.user", "SMTP_USER")
	viper.BindEnv("smtp.password", "SMTP_PASSWORD")

	viper.BindEnv("imap.host", "IMAP_HOST")
	viper.BindEnv("imap.port", "IMAP_PORT")
	viper.BindEnv("imap.user", "IMAP_USER")
	viper.BindEnv("imap.password", "IMAP_PASSWORD")

	viper.BindEnv("jwt.access_name", "JWT_ACCESS_NAME")
	viper.BindEnv("jwt.refresh_name", "JWT_REFRESH_NAME")
	viper.BindEnv("jwt.guest_name", "JWT_GUEST_NAME")
	viper.BindEnv("jwt.secret_key", "JWT_SECRET_KEY")
	viper.BindEnv("jwt.access_expires_in", "JWT_ACCESS_EXPIRES_IN")
	viper.BindEnv("jwt.refresh_expires_in", "JWT_REFRESH_EXPIRES_IN")

	viper.BindEnv("server.api_prefix", "SV_API_PREFIX")
	viper.BindEnv("server.port", "SV_PORT")
	viper.BindEnv("server.write_timeout", "SV_WRITE_TIMEOUT")
	viper.BindEnv("server.read_timeout", "SV_READ_TIMEOUT")
	viper.BindEnv("server.idle_timeout", "SV_IDLE_TIMEOUT")
	viper.BindEnv("server.allow_origins", "SV_ALLOW_ORIGINS")
	viper.BindEnv("server.allow_methods", "SV_ALLOW_METHODS")
	viper.BindEnv("server.allow_headers", "SV_ALLOW_HEADERS")
	viper.BindEnv("server.expose_headers", "SV_EXPOSE_HEADERS")
	viper.BindEnv("server.allow_credentials", "SV_ALLOW_CREDENTIALS")
	viper.BindEnv("server.max_age", "SV_MAX_AGE")
	viper.BindEnv("server.max_header_bytes", "SV_MAX_HEADER_BYTES")

	viper.AddConfigPath("./configs")
	viper.SetConfigName("main")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
