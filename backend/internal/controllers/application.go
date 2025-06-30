package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"home-server-hub/internal/models"
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
	apps.Get("/", c.listApplications)

	apps.Post("/", c.createApplication)

	// Outros endpoints serão implementados posteriormente
	// apps.Get("/:id", c.getApplication)
	// apps.Put("/:id", c.updateApplication)
	// apps.Delete("/:id", c.deleteApplication)
}

// discoverApplications descobre aplicações Docker
//	@Summary		Descobre aplicações Docker rodando no servidor
//	@Description	Busca por containers Docker em execução e retorna como aplicações potenciais
//	@Tags			applications
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.DiscoverResult
//	@Failure		500	{object}	map[string]string
//	@Router			/applications/discover [get]
func (c *ApplicationController) discoverApplications(ctx *fiber.Ctx) error {
	result, err := c.appService.DiscoverApplications()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Falha ao descobrir aplicações: " + err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(result)
}

// createApplication cria uma nova aplicação a partir de um container Docker
//	@Summary		Cria uma nova aplicação
//	@Description	Cria e armazena uma aplicação com base no ID de um container Docker e dados opcionais enviados
//	@Tags			applications
//	@Accept			json
//	@Produce		json
//	@Param			containerID	query		string					true	"ID do container Docker"
//	@Param			application	body		models.ApplicationInput	false	"Dados opcionais para sobrescrever valores padrão"
//	@Success		201			{object}	models.Application
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/applications [post]
func (c *ApplicationController) createApplication(ctx *fiber.Ctx) error {
	containerID := ctx.Query("container_id")
	if containerID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Parâmetro 'container_id' é obrigatório",
		})
	}

	var input models.ApplicationInput

	if err := ctx.BodyParser(input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido: " + err.Error(),
		})
	}

	application, err := c.appService.CreateApplication(containerID, &input)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao criar aplicação: " + err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(application)
}

// listApplications lista todas as aplicações cadastradas
//	@Summary		Lista aplicações
//	@Description	Retorna todas as aplicações já criadas e armazenadas no sistema
//	@Tags			applications
//	@Produce		json
//	@Success		200	{object}	models.ListApplicationsResult
//	@Failure		500	{object}	map[string]string
//	@Router			/applications [get]
func (c *ApplicationController) listApplications(ctx *fiber.Ctx) error {
	result, err := c.appService.ListApplications()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao listar aplicações: " + err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(result)
}
