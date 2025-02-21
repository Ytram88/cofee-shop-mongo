package config

import (
	"fmt"
	"os"
)

type Config struct {
	MongoUser     string
	MongoPassword string
	MongoHost     string
	MongoPort     string
	MongoDatabase string
}

var (
	cfg Config
)

func LoadConfig() *Config {
	cfg = Config{
		MongoUser:     getEnv("MONGO_USER", "cofeeStaff"),
		MongoPassword: getEnv("MONGO_PASSWORD", "pass123"),
		MongoHost:     getEnv("MONGO_HOST", "localhost"),
		MongoPort:     getEnv("MONGO_PORT", "27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "cofee-shop"),
	}
	return &cfg
}

func (c *Config) MakeConnectionString() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		c.MongoUser, c.MongoPassword, c.MongoHost, c.MongoPort, c.MongoDatabase)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
