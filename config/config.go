package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	RESTBindAddr      string
	ServerAppSOAPBase string
	MasterIP          string // IP para reemplazo en URLs (de .env)
}

func Load() Config {
	loadDotEnv(".env")

	// IP del Backend (desde .env)
	host := getEnv("SERVER_HOST", "127.0.0.1")
	port := getEnv("SERVER_PORT", "50021")

	// IP del servidor SOAP (desde .env)
	soapIP := getEnv("MASTER_IP", "127.0.0.1")

	return Config{
		RESTBindAddr:      "localhost" + ":" + port,
		ServerAppSOAPBase: strings.TrimRight(getEnv("SERVERAPP_SOAP_BASE", "http://"+soapIP+":8080/services"), "/"),
		MasterIP:          host, // Usamos la IP configurada para el reemplazo
	}
}

func loadDotEnv(path string) {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}

		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)

		// Expandir variables como ${MASTER_IP}
		value = os.ExpandEnv(value)

		_ = os.Setenv(key, value)
	}
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
