package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	RESTBindAddr      string
	ServerAppSOAPBase string
}

func Load() Config {
	loadDotEnv(".env")

	host := getEnv("SERVER_HOST", "127.0.0.1")
	port := getEnv("SERVER_PORT", "50021")

	return Config{
		RESTBindAddr:      host + ":" + port,
		ServerAppSOAPBase: strings.TrimRight(getEnv("SERVERAPP_SOAP_BASE", "http://localhost:8080/services"), "/"),
	}
}

func loadDotEnv(path string) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return
	}

	for _, rawLine := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" || os.Getenv(key) != "" {
			continue
		}

		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		_ = os.Setenv(key, value)
	}
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
