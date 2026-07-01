package core

import (
	"io"
	"olt-migrate-backend/internal/models"
)

// Parser define o contrato para extrair ONUs de um arquivo de backup
type Parser interface {
	Parse(file io.Reader) ([]models.ONU, error)
}

// Generator define o contrato para transformar ONUs no formato de CLI de destino
type Generator interface {
	Generate(onus []models.ONU) string
}
