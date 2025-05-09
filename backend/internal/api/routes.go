package api

import (
    "time"

	"github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/mongo" 

	"home-server-hub/internal/controllers"
	"home-server-hub/internal/docker"
	"home-server-hub/internal/repository"
	"home-server-hub/internal/services"
)

// setupApplicationRoutes configura as rotas relacionadas a aplicações
func setupApplicationRoutes(app *fiber.App, repo *repository.MongoApplicationRepository, dockerCli *docker.Client) {
	// Criar serviço de aplicação
	applicationService := services.NewApplicationService(repo, dockerCli)
	
	// Criar controller
	applicationController := controllers.NewApplicationController(applicationService)
	
	// Registrar rotas
	apiGroup := app.Group("/api/v1")
	applicationController.RegisterRoutes(apiGroup)
}

// SetupRoutes configura todas as rotas da API
func SetupRoutes(app *fiber.App, db *mongo.Database, dockerCli *docker.Client) {
	// Criar repositórios
	applicationRepo := repository.NewMongoApplicationRepository(db)
	
	// Configurar grupos de rotas
	setupApplicationRoutes(app, applicationRepo, dockerCli)
	
	// Adicionar rota de saúde
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}
