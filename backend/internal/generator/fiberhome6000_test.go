package generator

import (
	"strings"
	"testing"

	"olt-migrate-backend/internal/models"
)

func TestFiberhome6000Generator(t *testing.T) {
	onus := []models.ONU{
		{
			SlotID:       2,
			PortID:       3,
			OnuID:        15,
			SerialNumber: "FHTT9c0a4860",
			Model:        "HG6145E",
			PPPoEUser:    "maria.silva",
			PPPoEPass:    "123456", // Texto plano
			WiFiSSID:     "Maria_WiFi",
			WiFiPass:     "senha654",
			Services: []models.Service{
				{CVLAN: 189},
			},
		},
		{
			SlotID:       2,
			PortID:       3,
			OnuID:        16,
			SerialNumber: "FHTT11112222",
			Model:        "AN5506-01-A",
			// Sem PPPoE -> Bridge
		},
	}

	gen := &Fiberhome6000Generator{}
	output := gen.Generate(onus)

	// Verifica router AN6000
	if !strings.Contains(output, "whitelist add phy-id FHTT9c0a4860 checkcode fiberhome type HG6145E slot 2 pon 3 onuid 15") {
		t.Errorf("Falhou ao gerar whitelist AN6000: %s", output)
	}
	if !strings.Contains(output, "authorize 1/2/3 15 type HG6145E phy-id FHTT9c0a4860 password null") {
		t.Errorf("Falhou ao gerar authorize AN6000: %s", output)
	}
	if !strings.Contains(output, "dsp pppoe pro dis maria.silva 123456 null auto upnp_switch enable entries 6 fe1 fe2 fe3 fe4 ssid1 ssid5") {
		t.Errorf("Falhou ao gerar wancfg texto plano no gerador AN6000: %s", output)
	}
	if !strings.Contains(output, "ssid enable Maria_WiFi hide disable authmode wpa-psk/wpa2psk encrypt-type tkipaes wpakey senha654") {
		t.Errorf("Falhou ao gerar wifi no gerador AN6000: %s", output)
	}

	// Verifica bridge AN6000
	if !strings.Contains(output, "whitelist add phy-id FHTT11112222 type null slot 2 pon 3 onuid 16") {
		t.Errorf("Falhou ao gerar whitelist bridge AN6000: %s", output)
	}
	if !strings.Contains(output, "onu port vlan 16 eth 1 service 1 tag priority 255 tpid 33024 vid 10") {
		t.Errorf("Falhou ao gerar vlan bridge AN6000: %s", output)
	}
}
