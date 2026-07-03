package generator

import (
	"fmt"
	"sort"
	"strings"

	"olt-migrate-backend/internal/models"
)

// Fiberhome6000Generator gera comandos para OLTs linha AN6000, estruturados por seções limpas (Autorização, WAN, Wi-Fi, Gerenciamento) por PON
type Fiberhome6000Generator struct{}

// Generate recebe as ONUs mapeadas e retorna o script CLI no formato AN6000
func (g *Fiberhome6000Generator) Generate(onus []models.ONU) string {
	if len(onus) == 0 {
		return "! Nenhuma ONU para migrar.\n"
	}

	// 1. Ordenar determinística por Slot, Port e OnuID
	sortedONUs := make([]models.ONU, len(onus))
	copy(sortedONUs, onus)
	sort.Slice(sortedONUs, func(i, j int) bool {
		if sortedONUs[i].SlotID != sortedONUs[j].SlotID {
			return sortedONUs[i].SlotID < sortedONUs[j].SlotID
		}
		if sortedONUs[i].PortID != sortedONUs[j].PortID {
			return sortedONUs[i].PortID < sortedONUs[j].PortID
		}
		return sortedONUs[i].OnuID < sortedONUs[j].OnuID
	})

	// Agrupar ONUs por PON (Slot/Port)
	type ponKey struct {
		slot int
		port int
	}
	ponMap := make(map[ponKey][]models.ONU)
	var ponKeys []ponKey

	for _, onu := range sortedONUs {
		k := ponKey{slot: onu.SlotID, port: onu.PortID}
		if _, exists := ponMap[k]; !exists {
			ponKeys = append(ponKeys, k)
		}
		ponMap[k] = append(ponMap[k], onu)
	}

	var sb strings.Builder

	sb.WriteString("! ==============================================================================\n")
	sb.WriteString("! MIGRACAO PARA FIBERHOME AN6000 (AGRUPADO POR SECOES: AUTORIZACAO, WAN, GERENCIAMENTO)\n")
	sb.WriteString("! ==============================================================================\n\n")

	for _, k := range ponKeys {
		ponONUs := ponMap[k]
		slot := k.slot
		port := k.port

		sb.WriteString(fmt.Sprintf("! ==============================================================================\n"))
		sb.WriteString(fmt.Sprintf("! PORTA PON 1/%d/%d (%d ONUs)\n", slot, port, len(ponONUs)))
		sb.WriteString(fmt.Sprintf("! ==============================================================================\n\n"))

		// --- SEÇÃO 1: AUTORIZAÇÃO (Modo Global) ---
		sb.WriteString(fmt.Sprintf("! --- 1. AUTORIZACAO DAS ONUS (PON 1/%d/%d) ---\n", slot, port))
		for _, onu := range ponONUs {
			model := onu.Model
			if model == "" {
				model = "5506-04-FA"
			}
			serial := onu.SerialNumber
			if serial == "" {
				serial = "FHTT00000000"
			}

			if onu.PPPoEUser != "" {
				sb.WriteString(fmt.Sprintf("whitelist add phy-id %s checkcode fiberhome type %s slot %d pon %d onuid %d\n",
					serial, model, slot, port, onu.OnuID))
				sb.WriteString(fmt.Sprintf("authorize 1/%d/%d %d type %s phy-id %s password null\n",
					slot, port, onu.OnuID, model, serial))
				sb.WriteString(fmt.Sprintf("onu pon-type 1/%d/%d %d 712\n\n",
					slot, port, onu.OnuID))
			} else {
				sb.WriteString(fmt.Sprintf("whitelist add phy-id %s type null slot %d pon %d onuid %d\n",
					serial, slot, port, onu.OnuID))
				sb.WriteString(fmt.Sprintf("authorize 1/%d/%d %d type null phy-id %s password null\n",
					slot, port, onu.OnuID, serial))
				sb.WriteString(fmt.Sprintf("onu pon-type 1/%d/%d %d 712\n\n",
					slot, port, onu.OnuID))
			}
		}

		// --- SEÇÃO 2: WAN / INTERNET (Modo interface pon) ---
		sb.WriteString(fmt.Sprintf("! --- 2. CONFIGURACAO DE WAN / INTERNET (PON 1/%d/%d) ---\n", slot, port))
		sb.WriteString(fmt.Sprintf("interface pon 1/%d/%d\n\n", slot, port))
		for _, onu := range ponONUs {
			vlan := 188
			if onu.PPPoEUser == "" {
				vlan = 10
			}
			if len(onu.Services) > 0 && onu.Services[0].CVLAN > 0 {
				vlan = onu.Services[0].CVLAN
			}

			if onu.PPPoEUser != "" {
				login := onu.PPPoEUser
				senha := onu.PPPoEPass
				if senha == "" {
					senha = "123456"
				}
				sb.WriteString(fmt.Sprintf("! ONU %d (%s) - WAN\n", onu.OnuID, login))
				sb.WriteString(fmt.Sprintf("onu wan-cfg %d ind 1 mode tr069-int ty r %d 65535 nat en qos dis vlan tag tvlan dis 65535 65535 dsp pppoe pro dis %s %s null auto upnp_switch enable entries 6 fe1 fe2 fe3 fe4 ssid1 ssid5\n",
					onu.OnuID, vlan, login, senha))
				sb.WriteString(fmt.Sprintf("onu ipv6-wan-cfg %d ind 1 ip-stack-mode both ipv6-src-type dhcpv6 prefix-src-type delegate ipv6-address ::/0 ipv6-gateway :: ipv6-master-dns :: ipv6-slave-dns :: ipv6-static-prefix ::/0\n\n",
					onu.OnuID))
			} else {
				sb.WriteString(fmt.Sprintf("! ONU Bridge %d - LAN\n", onu.OnuID))
				sb.WriteString(fmt.Sprintf("onu port vlan %d eth 1 service count 1\n", onu.OnuID))
				sb.WriteString(fmt.Sprintf("onu port vlan %d eth 1 service 1 tag priority 255 tpid 33024 vid %d\n\n",
					onu.OnuID, vlan))
			}
		}

		// --- SEÇÃO 3: GERENCIAMENTO E ISOLAMENTO (Modo interface pon) ---
		sb.WriteString(fmt.Sprintf("! --- 3. GERENCIAMENTO E ISOLAMENTO (PON 1/%d/%d) ---\n", slot, port))
		for _, onu := range ponONUs {
			if onu.PPPoEUser != "" {
				sb.WriteString(fmt.Sprintf("! ONU %d (%s) - Gerenciamento\n", onu.OnuID, onu.PPPoEUser))
				sb.WriteString(fmt.Sprintf("onu port-isolation disable %d\n", onu.OnuID))
				sb.WriteString(fmt.Sprintf("onu remote-manage-cfg %d tr069 enable acs-url http://cwmp.nicnet.com.br:8088 acl-user Admin acl-pswd Admin@1234 inform enable interval 59834 port 0 user Admin pswd Admin@1234\n\n",
					onu.OnuID))
			}
		}

		sb.WriteString("exit\n\n")
	}

	return sb.String()
}
