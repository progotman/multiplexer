package provider

import (
	"os"
	"strconv"

	"github.com/progotman/multiplexer/config"
)

type EnvironmentProvider struct{}

func (c EnvironmentProvider) GetConfig() config.Config {
	return config.Config{
		ServerAddress:               getStringEnv("MULTIPLEXER_SERVER_ADDRESS"),
		MaxNumberOfUrls:             getIntegerEnv("MULTIPLEXER_MAX_NUMBER_OF_URLS"),
		MaxNumberOfIncomingRequests: getIntegerEnv("MULTIPLEXER_MAX_NUMBER_OF_INCOMING_REQUESTS"),
		MaxNumberOfOutgoingRequests: getIntegerEnv("MULTIPLEXER_MAX_NUMBER_OF_OUTGOING_REQUESTS"),
		RequestTimeout:              getIntegerEnv("MULTIPLEXER_REQUEST_TIMEOUT"),
	}
}

func getStringEnv(key string) string {
	return os.Getenv(key)
}

func getIntegerEnv(key string) int {
	val := getStringEnv(key)
	if value, err := strconv.Atoi(val); err == nil {
		return value
	}

	return 0
}
