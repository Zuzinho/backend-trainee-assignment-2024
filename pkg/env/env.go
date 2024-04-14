package env

import (
	"avito_hr/pkg/session"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func MustInit() {
	if err := godotenv.Load("./cmd/avito_hr/.env"); err != nil {
		log.Fatal(err)
	}
}

func MustJWTConfig() *session.JWTConfig {
	method, exist := os.LookupEnv("JWT_SIGNING_METHOD")
	if !exist {
		log.Fatal("no jwt signing method")
	}

	secret, exist := os.LookupEnv("JWT_SECRET_PASSWORD")
	if !exist {
		log.Fatal("no jwt secret password")
	}

	return session.NewJWTConfig(method, []byte(secret))
}

func MustDBConnectionString() string {
	conn, exist := os.LookupEnv("DB_CONNECTION_STRING")
	if !exist {
		log.Fatal("no db connection string")
	}

	return conn
}

func MustMaxConnectionCount() int {
	count, exist := os.LookupEnv("MAX_CONNECTIONS_COUNT")
	if !exist {
		log.Fatal("no max connection count")
	}

	d, err := strconv.Atoi(count)
	if err != nil {
		log.Fatal(err)
	}

	return d
}

func MustPort() string {
	port, exist := os.LookupEnv("PORT")
	if !exist {
		log.Fatal("no port")
	}

	return port
}
