package provider

import "github.com/progotman/multiplexer/config"

type InMemoryProvider struct{}

func (c InMemoryProvider) GetConfig() config.Config {
	return config.Config{
		ServerAddress:               ":8080",
		MaxNumberOfUrls:             20,
		MaxNumberOfIncomingRequests: 100,
		MaxNumberOfOutgoingRequests: 4,
		RequestTimeout:              1,
	}
}
