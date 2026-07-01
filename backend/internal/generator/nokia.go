package generator

import (
	"fmt"
	"strings"

	"olt-migrate-backend/internal/models"
)

type NokiaGenerator struct{}

// Generate recebe a lista de ONUs mapeadas e gera o script de CLI ISAM (Nokia).
func (g *NokiaGenerator) Generate(onus []models.ONU) string {
	var builder strings.Builder

	builder.WriteString("!========================================================\n")
	builder.WriteString("! SCRIPT GERADO AUTOMATICAMENTE: OLT MIGRATE (FIBERHOME -> NOKIA)\n")
	builder.WriteString("!========================================================\n\n")

	for _, onu := range onus {
		sn := onu.SerialNumber
		vendor := "FHTT"
		serialStr := sn
		
		// FHTT9df6baa0 -> Vendor: FHTT, Serial: 9df6baa0
		if len(sn) == 12 {
			vendor = sn[:4]
			serialStr = sn[4:]
		}

		s := onu.SlotID
		p := onu.PortID
		i := onu.OnuID
		
		// Fallback para VLAN (futuramente vindo do parser)
		vlan := 189

		// Determina se é Router (possui PPPoE/wan-cfg) ou Bridge
		isRouter := onu.PPPoEUser != ""

		builder.WriteString(fmt.Sprintf("! --- Provisionamento da ONU %d/%d/%d (SN: %s) ---\n", s, p, i, sn))
		
		if isRouter {
			// ==========================================
			// MODO ROUTER
			// ==========================================
			login := onu.PPPoEUser
			cliente := onu.PPPoEUser // Usando login como nome do cliente por enquanto
			
			builder.WriteString(fmt.Sprintf("configure equipment ont interface 1/1/%d/%d/%d sw-ver-pland auto desc1 \"%s\" desc2 \"%s\" sernum %s:%s sw-dnload-version auto voip-allowed veip pland-cfgfile1 DISABLED dnload-cfgfile1 DISABLED\n", s, p, i, login, cliente, vendor, serialStr))
			builder.WriteString(fmt.Sprintf("configure equipment ont interface 1/1/%d/%d/%d admin-state up\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure equipment ont slot 1/1/%d/%d/%d/1 planned-card-type ethernet plndnumdataports 3 plndnumvoiceports 0\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure equipment ont slot 1/1/%d/%d/%d/14 planned-card-type veip plndnumdataports 1 plndnumvoiceports 0 admin-state up\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure qos interface 1/1/%d/%d/%d/14/1 upstream-queue 0 bandwidth-profile name:HSI_1G_UP\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure interface port uni:1/1/%d/%d/%d/14/1 admin-up\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure bridge port 1/1/%d/%d/%d/14/1 max-unicast-mac 4 max-committed-mac 1\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure bridge port 1/1/%d/%d/%d/14/1 vlan-id %d tag single-tagged\n", s, p, i, vlan))
			
		} else {
			// ==========================================
			// MODO BRIDGE
			// ==========================================
			upProfile := "HSI_1G_UP"
			downProfile := "HSI_1G_DOWN"
			
			builder.WriteString(fmt.Sprintf("configure equipment ont interface 1/1/%d/%d/%d sw-ver-pland disabled sernum %s:%s\n", s, p, i, vendor, serialStr))
			builder.WriteString(fmt.Sprintf("configure equipment ont interface 1/1/%d/%d/%d admin-state up\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure equipment ont slot 1/1/%d/%d/%d/1 planned-card-type ethernet plndnumdataports 1 plndnumvoiceports 0 admin-state up\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure interface port uni:1/1/%d/%d/%d/1/1 admin-up\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure qos interface 1/1/%d/%d/%d/1/1 upstream-queue 0 bandwidth-profile name:%s\n", s, p, i, upProfile))
			builder.WriteString(fmt.Sprintf("configure qos interface 1/1/%d/%d/%d/1/1 queue 0 shaper-profile name:%s\n", s, p, i, downProfile))
			builder.WriteString(fmt.Sprintf("configure bridge port 1/1/%d/%d/%d/1/1 max-unicast-mac 4 max-committed-mac 1\n", s, p, i))
			builder.WriteString(fmt.Sprintf("configure bridge port 1/1/%d/%d/%d/1/1 vlan-id %d\n", s, p, i, vlan))
			builder.WriteString(fmt.Sprintf("configure bridge port 1/1/%d/%d/%d/1/1 pvid %d\n", s, p, i, vlan))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}
