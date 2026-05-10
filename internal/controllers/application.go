package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"home-server-hub/internal/models"
	"home-server-hub/internal/services"
	"home-server-hub/internal/utils"
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
	apps.Get("/:id/image", c.getApplicationImage)

	apps.Get("/:id", c.getApplication)
	apps.Post("/", c.createApplication)
	apps.Post("/:id/start", c.startApplication)
	apps.Post("/:id/stop", c.stopApplication)
	apps.Post("/:id/restart", c.restartApplication)
	apps.Put("/:id", c.updateApplication)
	apps.Delete("/:id", c.deleteApplication)
}

// discoverApplications descobre aplicações Docker
//
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
//
//	@Summary		Cria uma nova aplicação
//	@Description	Cria e armazena uma aplicação com base no ID de um container Docker e dados opcionais enviados via multipart/form-data
//	@Tags			applications
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			container_id	query		string	true	"ID do container Docker"
//	@Param			name			formData	string	false	"Nome personalizado da aplicação"
//	@Param			port			formData	uint16	false	"Porta a ser exposta"
//	@Param			url				formData	string	false	"URL pública da aplicação"
//	@Param			tags			formData	[]string	false	"Tags da aplicação (múltiplas ocorrências do campo)"	collectionFormat(multi)
//	@Param			image			formData	file	false	"Imagem opcional para representar a aplicação"
//	@Success		201				{object}	models.Application
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/applications/ [post]
func (c *ApplicationController) createApplication(ctx *fiber.Ctx) error {
	containerID := ctx.Query("container_id")
	if containerID == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Parâmetro 'container_id' é obrigatório",
		})
	}

	// Extrai campos do formulário
	var input models.ApplicationInput

	if name := ctx.FormValue("name"); name != "" {
		input.Name = &name
	}

	if portStr := ctx.FormValue("port"); portStr != "" {
		if portParsed, err := strconv.ParseUint(portStr, 10, 16); err == nil {
			port := uint16(portParsed)
			input.Port = &port
		}
	}

	if url := ctx.FormValue("url"); url != "" {
		input.URL = &url
	}

	input.Tags = parseTagsFromForm(ctx)

	// Extrai o arquivo da imagem
	fileHeader, _ := ctx.FormFile("image")

	imageData, err := utils.ParseImageFromFormFile(fileHeader)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao processar imagem: " + err.Error(),
		})
	}

	application, err := c.appService.CreateApplication(containerID, &input, imageData)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao criar aplicação: " + err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(application)
}

// getApplicationImage retorna o arquivo de imagem associado a uma aplicação
//
//	@Summary		Retorna a imagem de uma aplicação
//	@Tags			applications
//	@Produce		octet-stream
//	@Param			id	path		string	true	"ID da aplicação"
//	@Success		200	{file}		file
//	@Failure		404	{object}	map[string]string
//	@Router			/applications/{id}/image [get]
func (c *ApplicationController) getApplicationImage(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	path := c.appService.GetApplicationImagePath(id)
	if err := ctx.SendFile(path); err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Imagem não encontrada",
		})
	}
	return nil
}

// getApplication retorna uma aplicação específica
//
//	@Summary		Retorna uma aplicação
//	@Tags			applications
//	@Produce		json
//	@Param			id	path		string	true	"ID da aplicação"
//	@Success		200	{object}	models.Application
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/applications/{id} [get]
func (c *ApplicationController) getApplication(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	app, err := c.appService.GetApplication(id)
	if err != nil {
		if errors.Is(err, models.ErrApplicationNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Aplicação não encontrada",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao buscar aplicação: " + err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(app)
}

// updateApplication atualiza parcialmente uma aplicação existente
//
//	@Summary		Atualiza uma aplicação
//	@Description	Aceita os mesmos campos do POST via multipart/form-data; apenas os campos enviados são atualizados.
//	@Tags			applications
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string	true	"ID da aplicação"
//	@Param			name	formData	string	false	"Novo nome"
//	@Param			port	formData	uint16	false	"Nova porta"
//	@Param			url		formData	string	false	"Nova URL"
//	@Param			tags	formData	[]string	false	"Novas tags (substitui as existentes; envie sem valor para limpar)"	collectionFormat(multi)
//	@Param			image	formData	file	false	"Nova imagem"
//	@Success		200		{object}	models.Application
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/applications/{id} [put]
func (c *ApplicationController) updateApplication(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var input models.ApplicationInput

	if name := ctx.FormValue("name"); name != "" {
		input.Name = &name
	}

	if portStr := ctx.FormValue("port"); portStr != "" {
		if portParsed, err := strconv.ParseUint(portStr, 10, 16); err == nil {
			port := uint16(portParsed)
			input.Port = &port
		}
	}

	if url := ctx.FormValue("url"); url != "" {
		input.URL = &url
	}

	input.Tags = parseTagsFromForm(ctx)

	fileHeader, _ := ctx.FormFile("image")
	imageData, err := utils.ParseImageFromFormFile(fileHeader)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Erro ao processar imagem: " + err.Error(),
		})
	}

	app, err := c.appService.UpdateApplication(id, &input, imageData)
	if err != nil {
		if errors.Is(err, models.ErrApplicationNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Aplicação não encontrada",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao atualizar aplicação: " + err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(app)
}

// startApplication inicia o container associado a uma aplicação
//
//	@Summary	Inicia o container de uma aplicação
//	@Tags		applications
//	@Produce	json
//	@Param		id	path		string	true	"ID da aplicação"
//	@Success	200	{object}	models.Application
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/applications/{id}/start [post]
func (c *ApplicationController) startApplication(ctx *fiber.Ctx) error {
	return c.respondToContainerAction(ctx, c.appService.StartApplication, "iniciar")
}

// stopApplication para o container associado a uma aplicação
//
//	@Summary	Para o container de uma aplicação
//	@Tags		applications
//	@Produce	json
//	@Param		id	path		string	true	"ID da aplicação"
//	@Success	200	{object}	models.Application
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/applications/{id}/stop [post]
func (c *ApplicationController) stopApplication(ctx *fiber.Ctx) error {
	return c.respondToContainerAction(ctx, c.appService.StopApplication, "parar")
}

// restartApplication reinicia o container associado a uma aplicação
//
//	@Summary	Reinicia o container de uma aplicação
//	@Tags		applications
//	@Produce	json
//	@Param		id	path		string	true	"ID da aplicação"
//	@Success	200	{object}	models.Application
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/applications/{id}/restart [post]
func (c *ApplicationController) restartApplication(ctx *fiber.Ctx) error {
	return c.respondToContainerAction(ctx, c.appService.RestartApplication, "reiniciar")
}

func (c *ApplicationController) respondToContainerAction(
	ctx *fiber.Ctx,
	action func(string) (*models.Application, error),
	verb string,
) error {
	id := ctx.Params("id")
	app, err := action(id)
	if err != nil {
		if errors.Is(err, models.ErrApplicationNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Aplicação não encontrada",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao " + verb + " aplicação: " + err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(app)
}

// deleteApplication remove uma aplicação
//
//	@Summary		Remove uma aplicação
//	@Tags			applications
//	@Param			id	path	string	true	"ID da aplicação"
//	@Success		204
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/applications/{id} [delete]
func (c *ApplicationController) deleteApplication(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.appService.DeleteApplication(id); err != nil {
		if errors.Is(err, models.ErrApplicationNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Aplicação não encontrada",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Erro ao remover aplicação: " + err.Error(),
		})
	}
	return ctx.SendStatus(http.StatusNoContent)
}

// listApplications lista todas as aplicações cadastradas
//
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

// parseTagsFromForm extrai o campo `tags` de um multipart/form-data, suportando
// múltiplas ocorrências (`tags=foo&tags=bar`). Retorna nil se o campo não estiver
// presente; caso contrário, retorna a lista limpa (trim, sem vazios, deduplicada,
// preservando ordem). Slice vazio explícito sinaliza "limpar todas as tags".
func parseTagsFromForm(ctx *fiber.Ctx) []string {
	form, err := ctx.MultipartForm()
	if err != nil || form == nil {
		return nil
	}
	raw, ok := form.Value["tags"]
	if !ok {
		return nil
	}
	seen := make(map[string]struct{}, len(raw))
	cleaned := make([]string, 0, len(raw))
	for _, t := range raw {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, dup := seen[t]; dup {
			continue
		}
		seen[t] = struct{}{}
		cleaned = append(cleaned, t)
	}
	return cleaned
}
