package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GoServer struct {
		Host string
		Port string
	}
	PyServer struct {
		Host string
		Port string
	}
	WebHook struct {
		PayloadMaxSize int
		Burst          int
		RateLimit      int
		Secret         string
	}
}

func LoadEnv() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		slog.Info("Couldn't initialize godotenv. Skipping the loading.....")
	}

	var config Config

	config.GoServer.Host = os.Getenv("GO_SERVER_HOST")
	config.GoServer.Port = os.Getenv("GO_SERVER_PORT")

	config.PyServer.Host = os.Getenv("PYTHON_SERVER_HOST")
	config.PyServer.Port = os.Getenv("PYTHON_SERVER_PORT")

	var convErr error

	config.WebHook.Burst, convErr = strconv.Atoi(os.Getenv("WEBHOOK_BURST_RATE"))
	config.WebHook.PayloadMaxSize, convErr = strconv.Atoi(os.Getenv("WEBHOOK_MAX_PAYLOAD_SIZE"))
	config.WebHook.Secret = os.Getenv("GITHUB_WEBHOOK_SECRET")
	config.WebHook.RateLimit, convErr = strconv.Atoi(os.Getenv("WEBHOOK_RATE_LIMIT"))

	if convErr != nil {
		return nil, convErr
	}
	return &config, nil
}
