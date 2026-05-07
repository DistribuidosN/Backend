package utils

import (
	"regexp"
	"strings"
)

var ipRegex = regexp.MustCompile(`(https?://)(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|localhost|127\.0\.0\.1)(:\d+)?`)

// FixIP reemplaza cualquier IP/localhost en una URL por la MasterIP especificada,
// preservando el protocolo y el puerto, SOLO si la URL contiene rutas específicas.
func FixIP(url, masterIP string) string {
	if url == "" || masterIP == "" {
		return url
	}

	// Solo aplicar si la URL contiene las rutas de imágenes o exportaciones de MinIO
	if !strings.Contains(url, "/enfok-images") && !strings.Contains(url, "/exports") {
		return url
	}
	
	return ipRegex.ReplaceAllStringFunc(url, func(match string) string {
		parts := ipRegex.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		protocol := parts[1]
		// ip := parts[2]
		port := ""
		if len(parts) > 3 {
			port = parts[3]
		}
		
		return protocol + masterIP + port
	})
}

// FixIPInObject intenta corregir IPs en campos de texto de un objeto (usando reflexión o type assertion)
func FixIPInObject(obj interface{}, masterIP string) interface{} {
	if masterIP == "" {
		return obj
	}

	switch v := obj.(type) {
	case string:
		return FixIP(v, masterIP)
	case []interface{}:
		for i, item := range v {
			v[i] = FixIPInObject(item, masterIP)
		}
	case map[string]interface{}:
		for k, val := range v {
			v[k] = FixIPInObject(val, masterIP)
		}
	}
	// Nota: Para estructuras específicas, es mejor hacerlo manualmente en el handler
	// o usar reflexión profunda si se vuelve complejo.
	return obj
}
