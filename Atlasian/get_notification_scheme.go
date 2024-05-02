package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load("creds.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	// godotenv package
	dotenv := goDotEnvVariable("ATLASSIAN_API_KEY")

	fmt.Printf("godotenv : %s = %s \n", "ATLASSIAN_API_KEY", dotenv)
}
