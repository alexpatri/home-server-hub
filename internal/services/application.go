package services

import (
	"errors"
	"strings"

	"home-server-hub/internal/models"
	"home-server-hub/internal/utils/docker"
)

// ApplicationService encapsula a lógica de negócio relacionada a aplicações
type ApplicationService struct {
	repository models.ApplicationRepository
	dockerCli  *docker.Client
}

// NewApplicationService cria um novo serviço de aplicação
func NewApplicationService(repo models.ApplicationRepository, dockerCli *docker.Client) *ApplicationService {
	return &ApplicationService{
		repository: repo,
		dockerCli:  dockerCli,
	}
}

// DiscoverApplications descobre aplicações rodando como containers Docker
func (s *ApplicationService) DiscoverApplications() (*models.DiscoverResult, error) {
	// Obter containers ativos do Docker
	containers, err := s.dockerCli.ListContainers()
	if err != nil {
		return nil, err
	}

	// Obter containers já registrados no banco de dados
	existingContainers, err := s.repository.FindExistingContainers()
	if err != nil {
		return nil, err
	}

	// Processar e retornar resultado da descoberta
	result := &models.DiscoverResult{
		Discovered: make([]models.DiscoveredApplication, 0),
		Total:      len(containers),
	}

	for _, container := range containers {
		// Verifica se o container deve ser processado
		if !isSystemContainer(container.Name, container.Image) && !existingContainers[container.ID] {
			// Determinar porta principal
			var port uint16
			if len(container.Ports) > 0 {
				port = container.Ports[0].HostPort
			}

			// Criar tags baseadas no nome e imagem
			tags := generateTags(container.Name, container.Image)

			result.Discovered = append(result.Discovered, models.DiscoveredApplication{
				Name:      formatContainerName(container.Name),
				Container: container.ID,
				IP:        container.IP,
				Port:      port,
				Tags:      tags,
			})
		}
	}

	return result, nil
}

func (s *ApplicationService) CreateApplication(containerID string, input *models.ApplicationInput, image *models.Image) (*models.Application, error) {
	container, err := s.dockerCli.GetContainer(containerID)
	if err != nil {
		return nil, err
	}

	var defaultPort uint16
	if len(container.Ports) > 0 {
		defaultPort = container.Ports[0].HostPort
	}

	defaultName := formatContainerName(container.Name)

	application := &models.Application{
		Name:      deref(input.Name, defaultName),
		Container: container.ID,
		Port:      deref(input.Port, defaultPort),
		IP:        container.IP,
		Image:     image,
		URL:       deref(input.URL, ""),
		Tags:      input.Tags,
	}

	if err := s.repository.Create(application); err != nil {
		return nil, err
	}

	return application, nil
}

// GetApplicationImagePath retorna o caminho do arquivo de imagem de uma aplicação.
func (s *ApplicationService) GetApplicationImagePath(id string) string {
	return s.repository.ImagePath(id)
}

// GetApplication retorna uma aplicação pelo ID com o status calculado em tempo real.
func (s *ApplicationService) GetApplication(id string) (*models.Application, error) {
	app, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	app.Status, _ = s.dockerCli.GetContainerStatus(app.Container)
	return app, nil
}

// UpdateApplication aplica atualizações parciais a uma aplicação.
// Apenas campos não-nil em input são sobrescritos. Se image != nil, substitui
// a imagem atual; caso contrário a imagem existente é preservada.
func (s *ApplicationService) UpdateApplication(id string, input *models.ApplicationInput, image *models.Image) (*models.Application, error) {
	app, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		app.Name = *input.Name
	}
	if input.Port != nil {
		app.Port = *input.Port
	}
	if input.URL != nil {
		app.URL = *input.URL
	}
	if input.Tags != nil {
		app.Tags = input.Tags
	}
	if image != nil {
		app.Image = image
	}

	if err := s.repository.Update(app); err != nil {
		return nil, err
	}

	app.Status, _ = s.dockerCli.GetContainerStatus(app.Container)
	return app, nil
}

// DeleteApplication remove uma aplicação e seu arquivo de imagem.
func (s *ApplicationService) DeleteApplication(id string) error {
	return s.repository.Delete(id)
}

// StartApplication inicia o container associado a uma aplicação.
func (s *ApplicationService) StartApplication(id string) (*models.Application, error) {
	return s.applyContainerAction(id, s.dockerCli.StartContainer)
}

// StopApplication para o container associado a uma aplicação.
func (s *ApplicationService) StopApplication(id string) (*models.Application, error) {
	return s.applyContainerAction(id, s.dockerCli.StopContainer)
}

// RestartApplication reinicia o container associado a uma aplicação.
func (s *ApplicationService) RestartApplication(id string) (*models.Application, error) {
	return s.applyContainerAction(id, s.dockerCli.RestartContainer)
}

func (s *ApplicationService) applyContainerAction(id string, action func(string) error) (*models.Application, error) {
	app, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if app.Container == "" {
		return nil, errors.New("aplicação não possui container associado")
	}
	if err := action(app.Container); err != nil {
		return nil, err
	}
	app.Status, _ = s.dockerCli.GetContainerStatus(app.Container)
	return app, nil
}

func (s *ApplicationService) ListApplications() (*models.ListApplicationsResult, error) {
	applications, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	for i := range applications {
		applications[i].Status, _ = s.dockerCli.GetContainerStatus(applications[i].Container)
	}

	return &models.ListApplicationsResult{
		Applications: applications,
		Total:        len(applications),
	}, nil
}

// isSystemContainer verifica se o container é um container de sistema
func isSystemContainer(name, image string) bool {
	// Lista de nomes/imagens a serem ignorados
	systemNames := []string{
		"mongo", "mongodb", "postgres", "mysql", "mariadb",
		"influxdb", "prometheus", "grafana", "traefik", "nginx",
		"docker-proxy", "portainer", "watchtower", "home-server-hub",
	}

	nameLower := strings.ToLower(name)
	imageLower := strings.ToLower(image)

	for _, sys := range systemNames {
		if strings.Contains(nameLower, sys) || strings.Contains(imageLower, sys) {
			return true
		}
	}

	return false
}

// formatContainerName formata o nome do container para exibição
func formatContainerName(name string) string {
	// Remove underscores e hífens e capitaliza palavras
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-'
	})

	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[0:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}

// generateTags gera tags baseadas no nome e imagem do container
func generateTags(name, image string) []string {
	tags := make(map[string]bool)

	// Extrair nome da imagem sem a tag
	imageParts := strings.Split(image, ":")
	if len(imageParts) > 0 {
		imageBase := imageParts[0]
		imageName := imageBase[strings.LastIndex(imageBase, "/")+1:]
		if imageName != "" {
			tags[strings.ToLower(imageName)] = true
		}
	}

	// Adicionar tags baseadas no nome
	nameParts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-' || r == '.'
	})

	for _, part := range nameParts {
		if len(part) > 2 { // Ignora partes muito curtas
			tags[strings.ToLower(part)] = true
		}
	}

	// Converter mapa para slice
	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}

	return result
}

func deref[T any](ptr *T, fallback T) T {
	if ptr == nil {
		return fallback
	}
	return *ptr
}
