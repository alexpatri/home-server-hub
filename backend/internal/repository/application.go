package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"


	"home-server-hub/internal/models"
)

// MongoApplicationRepository implementa ApplicationRepository usando MongoDB
type MongoApplicationRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoApplicationRepository cria um novo repositório de aplicações
func NewMongoApplicationRepository(db *mongo.Database) *MongoApplicationRepository {
	return &MongoApplicationRepository{
	    db:         db,
		collection: db.Collection("applications"),
	}
}

// FindAll retorna todas as aplicações
func (r *MongoApplicationRepository) FindAll() ([]models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var applications []models.Application
	if err = cursor.All(ctx, &applications); err != nil {
		return nil, err
	}

	return applications, nil
}

// FindByID busca uma aplicação pelo ID
func (r *MongoApplicationRepository) FindByID(id primitive.ObjectID) (*models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var application models.Application
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&application)
	if err != nil {
		return nil, err
	}

	return &application, nil
}

// FindByContainer busca uma aplicação pelo ID do container
func (r *MongoApplicationRepository) FindByContainer(containerID string) (*models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var application models.Application
	err := r.collection.FindOne(ctx, bson.M{"container": containerID}).Decode(&application)
	if err != nil {
		return nil, err
	}

	return &application, nil
}

// Create cria uma nova aplicação
func (r *MongoApplicationRepository) Create(application *models.Application) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if application.ID.IsZero() {
		application.ID = primitive.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, application)
	return err
}

// Update atualiza uma aplicação existente
func (r *MongoApplicationRepository) Update(application *models.Application) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": application.ID}, application)
	return err
}

// Delete remove uma aplicação
func (r *MongoApplicationRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// FindExistingContainers retorna um mapa de IDs de containers já cadastrados
func (r *MongoApplicationRepository) FindExistingContainers() (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Projeção para buscar apenas o campo container
	projection := bson.M{"container": 1}
	cursor, err := r.collection.Find(ctx, bson.M{}, &options.FindOptions{
		Projection: projection,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type containerDoc struct {
		Container string `bson:"container"`
	}

	existingContainers := make(map[string]bool)
	for cursor.Next(ctx) {
		var doc containerDoc
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		existingContainers[doc.Container] = true
	}

	return existingContainers, nil
}
