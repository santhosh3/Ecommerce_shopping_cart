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
	HostMail               string
	HostPassword           string
	SMTPHost               string
	SMTPPort               int64
}

type getDataFromEnv interface {
	getEnv() interface{}
}

type getStringStruct struct {
	key      string
	fallback string
}

type getIntStruct struct {
	key      string
	fallback int64
}

func (s getStringStruct) getEnv() interface{} {
	if value, ok := os.LookupEnv(s.key); ok {
		return value
	}
	return s.fallback
}

func (s getIntStruct) getEnv() interface{} {
	if value, ok := os.LookupEnv(s.key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return s.fallback
		}
		return i
	}
	return s.fallback
}

func getEnvProd(s getDataFromEnv) interface{} {
	return s.getEnv()
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
    
	DB := "postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable"
	host := getStringStruct{key: "PUBLIC_HOST", fallback: "http://localhost"}
	port := getStringStruct{key: "PORT", fallback: "3501"}
	JWTSecret := getStringStruct{key: "JWT_SECRET", fallback: "ECOM"}
	PostgresString := getStringStruct{key: "POSTGRES_SQL", fallback: DB}
	JWTExpirationInSeconds := getIntStruct{key: "JWT_EXPIRATION", fallback: 3600 * 24 * 7}
	HostMail := getStringStruct{key: "host_mail", fallback: "santhoshchinna109@outlook.com"}
	HostPassword := getStringStruct{key: "host_password", fallback: "Chinna@123"}
	SMTPHost := getStringStruct{key: "SMTP_host", fallback: "smtp.office365.com"}
	SMTPPort := getIntStruct{key: "SMTP_port", fallback: 587}

	return Config{
		PublicHost:             getEnvProd(host).(string),
		Port:                   getEnvProd(port).(string),
		JWTSecret:              getEnvProd(JWTSecret).(string),
		PostgresString:         getEnvProd(PostgresString).(string),
		JWTExpirationInSeconds: getEnvProd(JWTExpirationInSeconds).(int64),
		HostMail:               getEnvProd(HostMail).(string),
		HostPassword:           getEnvProd(HostPassword).(string),
		SMTPHost:               getEnvProd(SMTPHost).(string),
		SMTPPort:               getEnvProd(SMTPPort).(int64),
	}
}
