package config

type Config struct {
	ServerAddress               string
	MaxNumberOfUrls             int
	MaxNumberOfIncomingRequests int
	MaxNumberOfOutgoingRequests int
	RequestTimeout              int
}

type Provider interface {
	GetConfig() Config
}
