package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Application representa uma aplicação no home server
type Application struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Tags      []string           `json:"tags" bson:"tags"`
	Image     *Image             `json:"image" bson:"image"`
	Container string             `json:"_" bson:"container"`
	IP        string             `json:"ip" bson:"ip"`
	Port      uint16             `json:"port" bson:"port"`
	URL       string             `json:"url" bson:"url"`
	Status    string             `json:"status" bson:"-"` // Status não é armazenado, é calculado em tempo real
}

type ApplicationInput struct {
	Name  *string  `json:"name,omitempty"`
	Port  *uint16  `json:"port,omitempty"`
	URL   *string  `json:"url,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

// Image representa a imagem/ícone associado à aplicação
type Image struct {
	Name   string `json:"name" bson:"name"`
	Data   []byte `json:"data" bson:"data"`
	Height int    `json:"height" bson:"height"`
	Width  int    `json:"width" bson:"width"`
	Size   int    `json:"size" bson:"size"`
}

// ApplicationRepository define a interface para acesso aos dados de aplicações
type ApplicationRepository interface {
	FindAll() ([]Application, error)
	FindByID(id primitive.ObjectID) (*Application, error)
	FindByContainer(containerID string) (*Application, error)
	Create(application *Application) error
	Update(application *Application) error
	Delete(id primitive.ObjectID) error
	FindExistingContainers() (map[string]bool, error)
}

// DiscoveredApplication representa uma aplicação descoberta pelo sistema
type DiscoveredApplication struct {
	Name      string   `json:"name"`
	Container string   `json:"container"`
	IP        string   `json:"ip"`
	Port      uint16   `json:"port"`
	Tags      []string `json:"tags"`
	Exists    bool     `json:"exists"`
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
