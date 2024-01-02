package utils

import "github.com/joho/godotenv"

// LoadEnv loads environment variables from the .env file
func LoadEnv() error {
	if err := godotenv.Load("../.env"); err != nil {
		return err
	}
	return nil
}
