package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort string
	DBConfig   DatabaseConfig
	MongoConfig MongoConfig
	JWTConfig  JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type MongoConfig struct {
	URI      string
	Database string
}

type JWTConfig struct {
	Secret        string
	RefreshSecret string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: GetEnv("PORT", "8080"),
		DBConfig: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			User:     GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", "postgres"),
			DBName:   GetEnv("DB_NAME", "prestasi_mahasiswa"),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
		MongoConfig: MongoConfig{
			URI:      GetEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: GetEnv("MONGO_DATABASE", "prestasi_mahasiswa"),
		},
		JWTConfig: JWTConfig{
			Secret:        GetEnv("JWT_SECRET", "your-secret-key"),
			RefreshSecret: GetEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
		},
	}
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBConfig.Host,
		c.DBConfig.Port,
		c.DBConfig.User,
		c.DBConfig.Password,
		c.DBConfig.DBName,
		c.DBConfig.SSLMode,
	)
}

func (c *Config) GetServerAddress() string {
	return ":" + c.ServerPort
}