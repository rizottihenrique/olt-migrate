package core

import (
	"olt-migrate-backend/internal/generator"
	"olt-migrate-backend/internal/parser"
)

// GetParser retorna o parser correto de acordo com a fabricante de origem
func GetParser(vendor string) Parser {
	switch vendor {
	case "Fiberhome", "Fiberhome AN5000", "Fiberhome AN6000":
		return &parser.FiberhomeParser{}
	default:
		return &parser.FiberhomeParser{}
	}
}

// GetGenerator retorna o gerador correto de acordo com a fabricante de destino
func GetGenerator(vendor string) Generator {
	switch vendor {
	case "Nokia":
		return &generator.NokiaGenerator{}
	case "Fiberhome AN5000":
		return &generator.Fiberhome5000Generator{}
	case "Fiberhome AN6000":
		return &generator.Fiberhome6000Generator{}
	// Futuramente:
	// case "Huawei":
	// 	return &generator.HuaweiGenerator{}
	default:
		// Default para Nokia
		return &generator.NokiaGenerator{}
	}
}
