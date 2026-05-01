package config

import (
	"os"
	"strings"
)

// Config armazena as configurações da aplicação
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Docker   DockerConfig
}

// ServerConfig armazena as configurações do servidor HTTP
type ServerConfig struct {
	Port string
}

// DatabaseConfig armazena as configurações do banco de dados
type DatabaseConfig struct {
	URI          string
	DatabaseName string
}

// DockerConfig armazena as configurações do cliente Docker
type DockerConfig struct {
	Host string
}

// LoadConfig carrega as configurações do ambiente
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8000"),
		},
		Database: DatabaseConfig{
			URI:          getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			DatabaseName: getEnv("MONGODB_DATABASE", "home_server_hub"),
		},
		Docker: DockerConfig{
			Host: getEnv("DOCKER_HOST", "unix:///var/run/docker.sock"),
		},
	}
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}
