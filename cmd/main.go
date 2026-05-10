package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"

	"home-server-hub/internal/api"
	"home-server-hub/internal/config"
	"home-server-hub/internal/database"
	"home-server-hub/internal/docker"
)

// @title			Home Server Hub API
// @version		0.1
// @description	Documentação da API Home Server Hub
// @host			localhost:8000
// @BasePath		/api/v1
func main() {
	// Flag para carregar .env
	useDotEnv := flag.Bool("dotenv", false, "Carregar variáveis do arquivo .env")
	flag.Parse()

	// Se a flag for true, carrega o .env
	if *useDotEnv {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Carregar configurações
	cfg := config.LoadConfig()

	// Garantir diretório de dados (banco e imagens)
	if err := os.MkdirAll(filepath.Dir(cfg.Database.Path), 0o755); err != nil {
		log.Fatalf("Erro ao criar diretório do banco: %v", err)
	}
	if err := os.MkdirAll(cfg.Database.ImagesDir, 0o755); err != nil {
		log.Fatalf("Erro ao criar diretório de imagens: %v", err)
	}

	// Conectar ao Banco de Dados
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	defer database.Disconnect(db)

	// Inicializar cliente Docker
	dockerCli, err := docker.NewClient(cfg.Docker.Host)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Docker: %v", err)
	}
	defer dockerCli.Close()

	// Criar e iniciar servidor
	server := api.NewServer(cfg)

	// Configurar rotas
	api.SetupRoutes(server.GetApp(), db, cfg.Database.ImagesDir, dockerCli)

	// Canal para receber sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Iniciar servidor em uma goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	log.Printf("Servidor iniciado na porta %s", cfg.Server.Port)

	// Aguardar sinal para encerrar
	<-quit
	log.Println("Encerrando servidor...")
}
