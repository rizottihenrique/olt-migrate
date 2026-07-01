package core

import (
	"olt-migrate-backend/internal/generator"
	"olt-migrate-backend/internal/parser"
)

// GetParser retorna o parser correto de acordo com a fabricante de origem
func GetParser(vendor string) Parser {
	switch vendor {
	case "Fiberhome":
		return &parser.FiberhomeParser{}
	// Futuramente:
	// case "ZTE":
	// 	return &parser.ZTEParser{}
	default:
		// Default para Fiberhome se não especificado ou não encontrado
		return &parser.FiberhomeParser{}
	}
}

// GetGenerator retorna o gerador correto de acordo com a fabricante de destino
func GetGenerator(vendor string) Generator {
	switch vendor {
	case "Nokia":
		return &generator.NokiaGenerator{}
	// Futuramente:
	// case "Huawei":
	// 	return &generator.HuaweiGenerator{}
	default:
		// Default para Nokia
		return &generator.NokiaGenerator{}
	}
}
