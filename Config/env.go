package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost                    string
	Port                          string
	RefreshJWTSecret              string
	AccessJWTSecret               string
	PostgresString                string
	HostMail                      string
	HostPassword                  string
	SMTPHost                      string
	SMTPPort                      int64
	AccessJWTExpirationInSeconds  int64
	RefreshJWTExpirationInSeconds int64
	RedisDB                       string
	GrpcPort                      string
	KafkaPort                     string
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

	//DB := "postgresql://postgres:postgres@db:5432/postgres?sslmode=disable"
	DB := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	RedisDB := getStringStruct{key: "REDIS_DB", fallback: "localhost:6379"}
	grpcPort := getStringStruct{key: "GRPC_PORT", fallback: "50051"}
	kafkaPort := getStringStruct{key: "KAFKA_PORT", fallback: "29092"}
	host := getStringStruct{key: "PUBLIC_HOST", fallback: "http://localhost"}
	port := getStringStruct{key: "PORT", fallback: "3500"}
	AccessJWTSecret := getStringStruct{key: "ACCESS_JWT_SECRET", fallback: "ACCESSTOKEN"}
	RefreshJWTSecret := getStringStruct{key: "REFRESH_JWT_SECRET", fallback: "REFRESHTOKEN"}
	PostgresString := getStringStruct{key: "POSTGRES_SQL", fallback: DB}
	HostMail := getStringStruct{key: "host_mail", fallback: "santhoshchinna109@outlook.com"}
	HostPassword := getStringStruct{key: "host_password", fallback: ""}
	SMTPHost := getStringStruct{key: "SMTP_host", fallback: "smtp.office365.com"}
	SMTPPort := getIntStruct{key: "SMTP_port", fallback: 587}
	AccessJWTExpirationInSeconds := getIntStruct{key: "ACCESS_JWT_EXPIRATION", fallback: 3600}
	RefreshJWTExpirationInSeconds := getIntStruct{key: "REFRESH_JWT_EXPIRATION", fallback: 3600 * 24 * 7}

	return Config{
		PublicHost:                    getEnvProd(host).(string),
		Port:                          getEnvProd(port).(string),
		AccessJWTSecret:               getEnvProd(AccessJWTSecret).(string),
		RefreshJWTSecret:              getEnvProd(RefreshJWTSecret).(string),
		PostgresString:                getEnvProd(PostgresString).(string),
		HostMail:                      getEnvProd(HostMail).(string),
		HostPassword:                  getEnvProd(HostPassword).(string),
		SMTPHost:                      getEnvProd(SMTPHost).(string),
		SMTPPort:                      getEnvProd(SMTPPort).(int64),
		RefreshJWTExpirationInSeconds: getEnvProd(RefreshJWTExpirationInSeconds).(int64),
		AccessJWTExpirationInSeconds:  getEnvProd(AccessJWTExpirationInSeconds).(int64),
		RedisDB:                       getEnvProd(RedisDB).(string),
		GrpcPort:                      getEnvProd(grpcPort).(string),
		KafkaPort:                     getEnvProd(kafkaPort).(string),
	}
}
