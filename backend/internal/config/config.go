package config

import "os"

type Config struct {
	APIAddr     string
	MetricsAddr string
}

func Load() Config {
	return Config{
		APIAddr:     getenv("CONDUCTOR_API_ADDR", ":5000"),
		MetricsAddr: getenv("CONDUCTOR_METRICS_ADDR", ":5000"),
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
