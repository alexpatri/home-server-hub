package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"home-server-hub/internal/models"
)

// SQLiteApplicationRepository implementa ApplicationRepository usando SQLite.
type SQLiteApplicationRepository struct {
	db        *sql.DB
	imagesDir string
}

// NewSQLiteApplicationRepository cria um novo repositório de aplicações.
func NewSQLiteApplicationRepository(db *sql.DB, imagesDir string) *SQLiteApplicationRepository {
	return &SQLiteApplicationRepository{db: db, imagesDir: imagesDir}
}

func (r *SQLiteApplicationRepository) ImagePath(id string) string {
	return filepath.Join(r.imagesDir, id)
}

func (r *SQLiteApplicationRepository) FindAll() ([]models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
        SELECT id, name, tags, container, ip, port, url,
               image_name, image_width, image_height, image_size
        FROM applications`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := make([]models.Application, 0)
	for rows.Next() {
		app, err := scanApplication(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, rows.Err()
}

func (r *SQLiteApplicationRepository) FindByID(id string) (*models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, tags, container, ip, port, url,
               image_name, image_width, image_height, image_size
        FROM applications WHERE id = ?`, id)

	app, err := scanApplication(row)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *SQLiteApplicationRepository) FindByContainer(containerID string) (*models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, `
        SELECT id, name, tags, container, ip, port, url,
               image_name, image_width, image_height, image_size
        FROM applications WHERE container = ?`, containerID)

	app, err := scanApplication(row)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *SQLiteApplicationRepository) Create(application *models.Application) error {
	if application.ID == "" {
		application.ID = uuid.NewString()
	}

	if err := r.persistImage(application); err != nil {
		return err
	}

	tagsJSON, err := marshalTags(application.Tags)
	if err != nil {
		return err
	}

	imageName, imgW, imgH, imgSize := imageMeta(application.Image)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = r.db.ExecContext(ctx, `
        INSERT INTO applications
            (id, name, tags, container, ip, port, url,
             image_name, image_width, image_height, image_size)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		application.ID, application.Name, tagsJSON, application.Container,
		application.IP, application.Port, application.URL,
		imageName, imgW, imgH, imgSize,
	)
	return err
}

func (r *SQLiteApplicationRepository) Update(application *models.Application) error {
	if application.ID == "" {
		return errors.New("application ID é obrigatório para update")
	}

	if err := r.persistImage(application); err != nil {
		return err
	}

	tagsJSON, err := marshalTags(application.Tags)
	if err != nil {
		return err
	}

	imageName, imgW, imgH, imgSize := imageMeta(application.Image)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = r.db.ExecContext(ctx, `
        UPDATE applications SET
            name = ?, tags = ?, container = ?, ip = ?, port = ?, url = ?,
            image_name = ?, image_width = ?, image_height = ?, image_size = ?
        WHERE id = ?`,
		application.Name, tagsJSON, application.Container,
		application.IP, application.Port, application.URL,
		imageName, imgW, imgH, imgSize,
		application.ID,
	)
	return err
}

func (r *SQLiteApplicationRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := r.db.ExecContext(ctx, `DELETE FROM applications WHERE id = ?`, id); err != nil {
		return err
	}
	if err := os.Remove(r.ImagePath(id)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (r *SQLiteApplicationRepository) FindExistingContainers() (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `SELECT container FROM applications WHERE container != ''`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existing := make(map[string]bool)
	for rows.Next() {
		var container string
		if err := rows.Scan(&container); err != nil {
			return nil, err
		}
		existing[container] = true
	}
	return existing, rows.Err()
}

// persistImage grava o arquivo de imagem em disco (se houver bytes) e descarta
// o conteúdo binário do model — só os metadados continuam em memória.
func (r *SQLiteApplicationRepository) persistImage(app *models.Application) error {
	if app.Image == nil || len(app.Image.Data) == 0 {
		return nil
	}
	if err := os.MkdirAll(r.imagesDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(r.ImagePath(app.ID), app.Image.Data, 0o644); err != nil {
		return err
	}
	app.Image.Data = nil
	return nil
}

// rowScanner abstrai *sql.Row e *sql.Rows para reutilizar scanApplication.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanApplication(s rowScanner) (models.Application, error) {
	var (
		app       models.Application
		tagsJSON  string
		imageName string
		imgW      int
		imgH      int
		imgSize   int
	)
	err := s.Scan(
		&app.ID, &app.Name, &tagsJSON, &app.Container,
		&app.IP, &app.Port, &app.URL,
		&imageName, &imgW, &imgH, &imgSize,
	)
	if err != nil {
		return app, err
	}

	tags, err := unmarshalTags(tagsJSON)
	if err != nil {
		return app, err
	}
	app.Tags = tags

	if imageName != "" || imgSize > 0 {
		app.Image = &models.Image{
			Name:   imageName,
			Width:  imgW,
			Height: imgH,
			Size:   imgSize,
		}
	}
	return app, nil
}

func marshalTags(tags []string) (string, error) {
	if tags == nil {
		return "[]", nil
	}
	b, err := json.Marshal(tags)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func unmarshalTags(raw string) ([]string, error) {
	if raw == "" {
		return []string{}, nil
	}
	var tags []string
	if err := json.Unmarshal([]byte(raw), &tags); err != nil {
		return nil, err
	}
	if tags == nil {
		tags = []string{}
	}
	return tags, nil
}

func imageMeta(img *models.Image) (name string, width, height, size int) {
	if img == nil {
		return "", 0, 0, 0
	}
	return img.Name, img.Width, img.Height, img.Size
}
