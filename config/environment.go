package config

import "os"

func GetEnvironment() string {
	environment := "local"
	envVariable := os.Getenv("APP_ENV")

	if envVariable == "development" {
		environment = "development"
	} else if envVariable == "production" {
		environment = "production"
	}

	return environment
}
