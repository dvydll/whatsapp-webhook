package config

import (
	"os"
	"strconv"
)

// Config contiene la configuración de la aplicación.
type Config struct {
	// WhatsApp
	WhatsAppVerifyToken   string
	WhatsAppPhoneNumberID string
	WhatsAppAccessToken   string
	WhatsAppAPIURL        string

	// TTS
	TTSEndpoint string
	TTSAPIKey   string

	// Server
	ServerPort string
}

// Load carga la configuración desde variables de entorno.
func Load() *Config {
	return &Config{
		// WhatsApp - soporta múltiples nombres de variables
		WhatsAppVerifyToken:   getEnv("WHATSAPP_VERIFY_TOKEN", "my_verify_token"),
		WhatsAppPhoneNumberID: getEnv("WHATSAPP_PHONE_NUMBER_ID", getEnv("PHONE_NUMBER_ID", "")),
		WhatsAppAccessToken:   getEnv("WHATSAPP_ACCESS_TOKEN", getEnv("META_ACCESS_TOKEN", "")),
		WhatsAppAPIURL:        getEnv("WHATSAPP_API_URL", "https://graph.facebook.com"),

		// TTS
		TTSEndpoint: getEnv("TTS_ENDPOINT", "http://localhost:8080"),
		TTSAPIKey:   getEnv("TTS_API_KEY", ""),

		// Server
		ServerPort: getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
