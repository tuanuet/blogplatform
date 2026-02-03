package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Logger    LoggerConfig
	Telemetry TelemetryConfig
	Scheduler SchedulerConfig
	Firebase  FirebaseConfig
	SePay     SePayConfig
}

// SePayConfig holds SePay-related configuration
type SePayConfig struct {
	APIKey       string `mapstructure:"api_key"`
	WebhookToken string `mapstructure:"webhook_token"`
	BankName     string `mapstructure:"bank_name"`
	BankAccount  string `mapstructure:"bank_account"`
	BankOwner    string `mapstructure:"bank_owner"`
	BankBranch   string `mapstructure:"bank_branch"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"` // debug, release, test
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig holds redis-related configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LoggerConfig holds logger-related configuration
type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json, text
}

// TelemetryConfig holds telemetry-related configuration
type TelemetryConfig struct {
	ServiceName string `mapstructure:"service_name"`
	Enabled     bool   `mapstructure:"enabled"`
}

// SchedulerConfig holds scheduler-related configuration
type SchedulerConfig struct {
	Enabled                bool   `mapstructure:"enabled"`
	DailyRecalculationHour int    `mapstructure:"daily_recalculation_hour"`
	Timezone               string `mapstructure:"timezone"`
}

// FirebaseConfig holds Firebase Cloud Messaging configuration
type FirebaseConfig struct {
	ProjectID          string `mapstructure:"project_id"`
	Enabled            bool   `mapstructure:"enabled"`
	APIKey             string `mapstructure:"api_key"`              // For service account authentication
	ServiceAccountPath string `mapstructure:"service_account_path"` // Path to service account JSON
}

// Load loads the configuration from file and environment variables
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Read from environment variables
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("config file not found, using defaults: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "boilerplate")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Logger defaults
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.format", "json")

	// Telemetry defaults
	viper.SetDefault("telemetry.service_name", "go-boilerplate")
	viper.SetDefault("telemetry.enabled", true)

	// Scheduler defaults
	viper.SetDefault("scheduler.enabled", true)
	viper.SetDefault("scheduler.daily_recalculation_hour", 0)
	viper.SetDefault("scheduler.timezone", "Local")

	// Firebase defaults
	viper.SetDefault("firebase.enabled", false)
	viper.SetDefault("firebase.project_id", "")
	viper.SetDefault("firebase.api_key", "")
	viper.SetDefault("firebase.service_account_path", "")

	// SePay defaults
	viper.SetDefault("sepay.api_key", "")
	viper.SetDefault("sepay.webhook_token", "")
	viper.SetDefault("sepay.bank_name", "")
	viper.SetDefault("sepay.bank_account", "")
	viper.SetDefault("sepay.bank_owner", "")
	viper.SetDefault("sepay.bank_branch", "")
}
