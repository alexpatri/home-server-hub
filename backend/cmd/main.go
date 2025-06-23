package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"home-server-hub/internal/api"
	"home-server-hub/internal/config"
	"home-server-hub/internal/docker"
)

// @title           Home Server Hub API
// @version         0.1
// @description     Documentação da API Home Server Hub
// @host            localhost:8080
// @BasePath        /api/v1
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

	// Conectar ao MongoDB
	db, err := connectDB(cfg.Database)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer disconnectDB(db)

	// Inicializar cliente Docker
	dockerCli, err := docker.NewClient(cfg.Docker.Host)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Docker: %v", err)
	}
	defer dockerCli.Close()

	// Criar e iniciar servidor
	server := api.NewServer(cfg)

	// Configurar rotas
	api.SetupRoutes(server.GetApp(), db, dockerCli)

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

// connectDB estabelece conexão com o MongoDB
func connectDB(dbConfig config.DatabaseConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(dbConfig.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verificar conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Conectado ao MongoDB!")
	return client.Database(dbConfig.DatabaseName), nil
}

// disconnectDB encerra a conexão com o MongoDB
func disconnectDB(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Client().Disconnect(ctx); err != nil {
		log.Fatalf("Erro ao desconectar do MongoDB: %v", err)
	}
	fmt.Println("Desconectado do MongoDB")
}
