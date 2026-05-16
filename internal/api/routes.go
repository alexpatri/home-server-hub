package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "home-server-hub/docs"
	"home-server-hub/internal/events"
	"home-server-hub/internal/handlers"
	"home-server-hub/internal/repository"
	"home-server-hub/internal/services"
	"home-server-hub/internal/utils/docker"
)

// setupApplicationRoutes configura as rotas relacionadas a aplicações
func setupApplicationRoutes(app *fiber.App, repo *repository.SQLiteApplicationRepository, dockerCli *docker.Client, broadcaster *events.Broadcaster) {
	applicationService := services.NewApplicationService(repo, dockerCli)
	applicationController := handlers.NewApplicationController(applicationService, broadcaster)
	apiGroup := app.Group("/api/v1")
	applicationController.RegisterRoutes(apiGroup)
}

// SetupRoutes configura todas as rotas da API
func SetupRoutes(app *fiber.App, repo *repository.SQLiteApplicationRepository, dockerCli *docker.Client, broadcaster *events.Broadcaster) {
	setupApplicationRoutes(app, repo, dockerCli, broadcaster)

	// Rota de saúde
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Rota /docs que redireciona para a interface Swagger
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html", 301)
	})

	// Swagger handler
	app.Get("/docs/*", swagger.HandlerDefault)
}
