package main

import (
	"log"
	"os"
)

// Configuration -- Конфигурация приложеиния
type Configuration struct {
	Path                    string
	Port                    string
	VkClientID              string
	VkClientSecret          string
	FireBaseCredentialsFile string
}

// NewConfiguration -- Создает новую конфигурацию
func NewConfiguration() *Configuration {
	vkClientID := lookupEnvOrDefault("VK_CLIENT_ID", "")
	if len(vkClientID) == 0 {
		log.Panic("Parameter 'VK_CLIENT_ID' is not set or invalid")
	}

	vkClientSecret := lookupEnvOrDefault("VK_CLIENT_SECRET", "")
	if len(vkClientSecret) == 0 {
		log.Panic("Parameter 'VK_CLIENT_SECRET' is not set or invalid")
	}

	return &Configuration{
		Path:                    lookupEnvOrDefault("SERVER_PATH", "/"),
		Port:                    lookupEnvOrDefault("SERVER_PORT", ":8080"),
		VkClientID:              vkClientID,
		VkClientSecret:          vkClientSecret,
		FireBaseCredentialsFile: lookupEnvOrDefault("FIREBASE_CREDENTIALS_FILE", ""),
	}
}

func lookupEnvOrDefault(key string, defaultValue string) string {
	value, present := os.LookupEnv(key)
	if !present {
		value = defaultValue
	}
	return value
}
