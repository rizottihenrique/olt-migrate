package generator

import (
	"strings"
	"testing"
	"olt-migrate-backend/internal/models"
)

func TestGenerateNokiaCommands_Router(t *testing.T) {
	onus := []models.ONU{
		{
			SlotID:       3,
			PortID:       1,
			OnuID:        1,
			SerialNumber: "FHTT9df6baa0",
			PPPoEUser:    "user.teste",
			PPPoEPass:    "1234",
		},
	}

	generatorInstance := &NokiaGenerator{}
	result := generatorInstance.Generate(onus)

	if !strings.Contains(result, "interface 1/1/3/1/1") {
		t.Errorf("A interface ISAM Nokia não foi gerada corretamente: %s", result)
	}

	if !strings.Contains(result, "FHTT:9df6baa0") {
		t.Errorf("O formato do Serial Number não foi ajustado corretamente para Nokia: %s", result)
	}

	if !strings.Contains(result, "user.teste") {
		t.Errorf("O login do cliente não foi incluído no desc: %s", result)
	}
}

func TestGenerateNokiaCommands_Bridge(t *testing.T) {
	onus := []models.ONU{
		{
			SlotID:       3,
			PortID:       1,
			OnuID:        2,
			SerialNumber: "FHTT12345678",
			PPPoEUser:    "", // Vazio = Bridge
		},
	}

	generatorInstance := &NokiaGenerator{}
	result := generatorInstance.Generate(onus)

	if !strings.Contains(result, "interface 1/1/3/1/2") {
		t.Errorf("A interface ISAM Nokia não foi gerada corretamente: %s", result)
	}

	if strings.Contains(result, "desc1") {
		t.Errorf("Bridge não deveria ter desc1 no comando base: %s", result)
	}

	if !strings.Contains(result, "pvid 189") {
		t.Errorf("Bridge deveria configurar o pvid: %s", result)
	}
}
