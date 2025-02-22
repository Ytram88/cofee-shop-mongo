package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
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
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}
	cfg = Config{
		MongoUser:     getEnv("MONGO_USER", "cofeeStaff"),
		MongoPassword: getEnv("MONGO_PASSWORD", "pass123"),
	}
	return &cfg
}

func (c *Config) MakeConnectionString() string {
	return fmt.Sprintf("mongodb+srv://%s:%s@cluster0.whhpn.mongodb.net/",
		c.MongoUser, c.MongoPassword)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
