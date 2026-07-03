package generator

import (
	"strings"
	"testing"

	"olt-migrate-backend/internal/models"
)

func TestFiberhome5000Generator(t *testing.T) {
	onus := []models.ONU{
		{
			SlotID:       1,
			PortID:       1,
			OnuID:        2,
			SerialNumber: "FHTT9a980b98",
			Model:        "5506-04-FA",
			PPPoEUser:    "cliente.pppoe",
			PPPoEPass:    "senha123", // Texto plano
			WiFiSSID:     "Cliente_WiFI",
			WiFiPass:     "minhasenha",
			Services: []models.Service{
				{CVLAN: 81},
			},
		},
		{
			SlotID:       1,
			PortID:       2,
			OnuID:        10,
			SerialNumber: "FHTT92370718",
			Model:        "5506-01-A1",
			// Sem PPPoE -> Bridge
		},
	}

	gen := &Fiberhome5000Generator{}
	output := gen.Generate(onus)

	// Verifica router AN5000
	if !strings.Contains(output, "set whitelist phy_addr address FHTT9a980b98 password null action add slot 1 pon 1 onu null type null") {
		t.Errorf("Falhou ao gerar whitelist AN5000: %s", output)
	}
	if !strings.Contains(output, "dsp pppoe pro dis cliente.pppoe senha123 null auto entries 6 fe1 fe2 fe3 fe4 ssid1 ssid5") {
		t.Errorf("Falhou ao gerar wancfg com senha texto plano no gerador AN5000: %s", output)
	}
	if strings.Contains(output, "wifi_serv_wlan") {
		t.Errorf("Gerador AN5000 não deveria incluir comandos de Wi-Fi: %s", output)
	}

	// Verifica bridge AN5000
	if !strings.Contains(output, "set epon slot 1 pon 2 onu 10 port 1 service number 1") {
		t.Errorf("Falhou ao gerar comando bridge AN5000: %s", output)
	}
}
