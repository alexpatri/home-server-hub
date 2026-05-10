package models

import "errors"

// ErrApplicationNotFound é retornado quando uma aplicação não existe.
var ErrApplicationNotFound = errors.New("application not found")

// Application representa uma aplicação no home server
type Application struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Tags      []string `json:"tags"`
	Image     *Image   `json:"image"`
	Container string   `json:"-"`
	IP        string   `json:"ip"`
	Port      uint16   `json:"port"`
	URL       string   `json:"url"`
	Status    string   `json:"status"` // calculado em tempo real, não persistido
}

type ApplicationInput struct {
	Name *string  `json:"name,omitempty"`
	Port *uint16  `json:"port,omitempty"`
	URL  *string  `json:"url,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

// Image representa a imagem/ícone associado à aplicação
type Image struct {
	Name   string `json:"name"`
	Data   []byte `json:"-"` // só usado durante o upload; persistido em disco
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Size   int    `json:"size"`
}

// ApplicationRepository define a interface para acesso aos dados de aplicações
type ApplicationRepository interface {
	FindAll() ([]Application, error)
	FindByID(id string) (*Application, error)
	FindByContainer(containerID string) (*Application, error)
	Create(application *Application) error
	Update(application *Application) error
	Delete(id string) error
	FindExistingContainers() (map[string]bool, error)
	ImagePath(id string) string
}

// DiscoveredApplication representa uma aplicação descoberta pelo sistema
type DiscoveredApplication struct {
	Name      string   `json:"name"`
	Container string   `json:"container"`
	IP        string   `json:"ip"`
	Port      uint16   `json:"port"`
	Tags      []string `json:"tags"`
}

// DiscoverResult representa o resultado da operação de descoberta
type DiscoverResult struct {
	Discovered []DiscoveredApplication `json:"discovered"`
	Total      int                     `json:"total"`
}

type ListApplicationsResult struct {
	Applications []Application `json:"applications"`
	Total        int           `json:"total"`
}
