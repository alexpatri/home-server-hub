package utils

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"

	"home-server-hub/internal/models"
)

// ParseImageFromFormFile processa um arquivo de imagem enviado via multipart e retorna uma struct Image
func ParseImageFromFormFile(fileHeader *multipart.FileHeader) (*models.Image, error) {
	if fileHeader == nil {
		return nil, nil
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		cfg.Width = 0
		cfg.Height = 0
	}

	return &models.Image{
		Name:   fileHeader.Filename,
		Data:   data,
		Size:   int(fileHeader.Size),
		Width:  cfg.Width,
		Height: cfg.Height,
	}, nil
}
