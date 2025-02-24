package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type JWTConfig struct {
	JWTSecret              string
	JWTExpirationInSeconds int64
}

type Config struct {
	Host          string
	Port          string
	MongoUser     string
	MongoPassword string
	JWTConfig     JWTConfig
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}
	jwtcfg := JWTConfig{
		JWTSecret:              getEnv("JWT_SECRET", "secretnword123"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
	}
	cfg := Config{
		MongoUser:     getEnv("MONGO_USER", "cofeeStaff"),
		MongoPassword: getEnv("MONGO_PASSWORD", "pass123"),
		Port:          getEnv("PORT", "8080"),
		JWTConfig:     jwtcfg,
	}
	return &cfg
}

// MakeConnectionString Only for MongoAtlas, change the function dependent on the database you use
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

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
