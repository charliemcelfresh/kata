package config

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	requiredEnvVars = [2]string{
		"DATABASE_URL",
		"SERVER_PORT",
	}
	Constants = map[string]interface{}{
		"REQUIRED_API_KATA_REQUEST_CONTENT_TYPE": "application/api.kata-json",
		"API_KATA_RESPONSE_CONTENT_TYPE":         "application/api.kata-json",
	}
	DB *sql.DB
)

func init() {
	logrus.SetReportCaller(true)
	err := godotenv.Overload()
	if err != nil {
		logrus.Fatal(err)
	}
	loadEnvVars()
	loadConfig()
}

func loadEnvVars() {
	for _, varName := range requiredEnvVars {
		varValue := os.Getenv(varName)
		if varValue == "" {
			err := errors.New(fmt.Sprintf("ENV VAR %v does not exist", varName))
			logrus.Fatal(err)
		}
		Constants[varName] = varValue
	}
}

func loadConfig() {
	err := godotenv.Overload()
	if err != nil {
		logrus.Fatal("Cannot load .env")
	}
	DB, err = sql.Open("postgres", fmt.Sprintf("%v", Constants["DATABASE_URL"]))
	if err != nil {
		logrus.Fatal("error connecting to database")
	}
}
