package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"home-server-hub/internal/services"
)

// ApplicationController gerencia as rotas relacionadas a aplicações
type ApplicationController struct {
	appService *services.ApplicationService
}

// NewApplicationController cria um novo controller de aplicações
func NewApplicationController(appService *services.ApplicationService) *ApplicationController {
	return &ApplicationController{
		appService: appService,
	}
}

// RegisterRoutes registra as rotas do controller no router do Fiber
func (c *ApplicationController) RegisterRoutes(router fiber.Router) {
	apps := router.Group("/applications")
	apps.Get("/discover", c.discoverApplications)

	// Outros endpoints serão implementados posteriormente
	// apps.Get("/", c.getAllApplications)
	// apps.Post("/", c.createApplication)
	// apps.Get("/:id", c.getApplication)
	// apps.Put("/:id", c.updateApplication)
	// apps.Delete("/:id", c.deleteApplication)
}

// discoverApplications descobre aplicações Docker
// @Summary Descobre aplicações Docker rodando no servidor
// @Description Busca por containers Docker em execução e retorna como aplicações potenciais
// @Tags applications
// @Accept json
// @Produce json
// @Success 200 {object} models.DiscoverResult
// @Failure 500 {object} map[string]string
// @Router /applications/discover [post]
func (c *ApplicationController) discoverApplications(ctx *fiber.Ctx) error {
	result, err := c.appService.DiscoverApplications()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Falha ao descobrir aplicações: " + err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(result)
}
