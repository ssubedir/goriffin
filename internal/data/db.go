package data

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DBConnection() map[string]interface{} {

	Connection := make(map[string]interface{})

	errr := godotenv.Load()
	if errr != nil {
		log.Fatal("Error loading .env file")
	}

	mongoConn := Mongo{
		Username: os.Getenv("DB0_USER"),
		Password: os.Getenv("DB0_PASSWORD"),
		Cluster:  os.Getenv("DB0_CLUSTER"),
	}

	mongoConnection := mongoConn

	Connection["mongodb"] = mongoConnection

	return Connection
}
