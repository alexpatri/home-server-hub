package api

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "home-server-hub/docs"
	"home-server-hub/internal/controllers"
	"home-server-hub/internal/docker"
	"home-server-hub/internal/repository"
	"home-server-hub/internal/services"
)

// setupApplicationRoutes configura as rotas relacionadas a aplicações
func setupApplicationRoutes(app *fiber.App, repo *repository.SQLiteApplicationRepository, dockerCli *docker.Client) {
	applicationService := services.NewApplicationService(repo, dockerCli)
	applicationController := controllers.NewApplicationController(applicationService)
	apiGroup := app.Group("/api/v1")
	applicationController.RegisterRoutes(apiGroup)
}

// SetupRoutes configura todas as rotas da API
func SetupRoutes(app *fiber.App, db *sql.DB, imagesDir string, dockerCli *docker.Client) {
	applicationRepo := repository.NewSQLiteApplicationRepository(db, imagesDir)

	setupApplicationRoutes(app, applicationRepo, dockerCli)

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
