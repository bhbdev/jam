package config

import (
	"fmt"

	"github.com/bhbdev/jam/lib/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Server      Server
	Database    Database
	Development bool
}

func DefaultConfig() *Config {
	return &Config{
		Server: Server{
			Host: "0.0.0.0",
			Port: 9090,
		},
		Database: Database{
			Engine: "sqlite3",
			Name:   "jam.db",
		},
		Development: true,
	}
}

func Load() (*Config, error) {
	var cfg = DefaultConfig()

	// Load configuration from yaml file
	// and parse it into config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // read in environment variables that match

	viper.BindEnv("Server.Host", "SERVER_HOST")
	viper.BindEnv("Server.Port", "SERVER_PORT")
	viper.BindEnv("Database.Engine", "DB_ENGINE")
	viper.BindEnv("Database.Name", "DB_NAME")
	viper.BindEnv("Database.User", "DB_USER")
	viper.BindEnv("Database.Pass", "DB_PASS")
	viper.BindEnv("Database.Host", "DB_HOST")
	viper.BindEnv("Database.Port", "DB_PORT")
	viper.BindEnv("Development", "DEVELOPMENT")

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			logger.Info("No config file found. Using defaults and environment variables")
		case viper.ConfigParseError:
			logger.Error("Error parsing config file. Please make sure the config file is well formed")
		default:
			logger.Error("fatal error loading config file:", "error", err)
		}
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error("fatal error unmarshalling config file:", "error", err)
		return nil, err
	}

	return cfg, nil
}

type Server struct {
	Host string
	Port int
}

func (s Server) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type Database struct {
	Engine string
	Name   string
	User   string
	Pass   string
	Host   string
	Port   int
	// pg only
	SslMode string
}

func (d Database) DSN() string {
	switch d.Engine {
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?multiStatements=true&parseTime=true",
			d.User,
			d.Pass,
			d.Host,
			d.Port,
			d.Name,
		)
	case "postgres":
		return fmt.Sprintf(
			"%s:%s@%s:%d/%s?sslmode=%s",
			d.User,
			d.Pass,
			d.Host,
			d.Port,
			d.Name,
			d.SslMode,
		)
	case "sqlite3":
		return fmt.Sprintf("file:%s?_foreign_keys=true", d.Name)
	default:
		return ""
	}
}
