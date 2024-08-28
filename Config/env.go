package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	JWTSecret              string
	JWTExpirationInSeconds int64
	PostgresString         string
}

var Envs = initConfig()


func initConfig() Config {
	godotenv.Load();

	   //localPostgresDB := "postgresql://postgres:postgres@localhost:5435/ecom?sslmode=disable";
        dockerPostgresDB := "postgresql://postgres:postgres@db:5432/ecom?sslmode=disable";

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "3500"),
		JWTSecret:              getEnv("JWT_SECRET", "ECOM"),
		PostgresString:         getEnv("POSTGRES_SQL", dockerPostgresDB),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION", 3600*24*7),
	}
}


// get the ENV by key or fallback

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
