package provider

import "github.com/progotman/multiplexer/config"

type FallbackProvider struct {
	Provider         config.Provider
	FallbackProvider config.Provider
}

func (p FallbackProvider) GetConfig() config.Config {
	originalConfig := p.Provider.GetConfig()
	fallbackConfig := p.FallbackProvider.GetConfig()

	return config.Config{
		ServerAddress:               fallbackString(originalConfig.ServerAddress, fallbackConfig.ServerAddress),
		MaxNumberOfUrls:             fallbackInt(originalConfig.MaxNumberOfUrls, fallbackConfig.MaxNumberOfUrls),
		MaxNumberOfIncomingRequests: fallbackInt(originalConfig.MaxNumberOfIncomingRequests, fallbackConfig.MaxNumberOfIncomingRequests),
		MaxNumberOfOutgoingRequests: fallbackInt(originalConfig.MaxNumberOfOutgoingRequests, fallbackConfig.MaxNumberOfOutgoingRequests),
		RequestTimeout:              fallbackInt(originalConfig.RequestTimeout, fallbackConfig.RequestTimeout),
	}
}

func fallbackInt(originalValue int, fallbackValue int) int {
	if originalValue == 0 {
		return fallbackValue
	}

	return originalValue
}

func fallbackString(originalValue string, fallbackValue string) string {
	if originalValue == "" {
		return fallbackValue
	}

	return originalValue
}
