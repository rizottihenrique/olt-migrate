package generator

import (
	"fmt"
	"sort"
	"strings"

	"olt-migrate-backend/internal/models"
)

// Fiberhome5000Generator gera comandos para OLTs linha AN5000, estruturados por seções limpas (Autorização, WAN, Wi-Fi, Gerenciamento) por PON
type Fiberhome5000Generator struct{}

// Generate recebe as ONUs mapeadas e retorna o script CLI no formato AN5000
func (g *Fiberhome5000Generator) Generate(onus []models.ONU) string {
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
	sb.WriteString("! MIGRACAO PARA FIBERHOME AN5000 (AGRUPADO POR SECOES: AUTORIZACAO, WAN, GERENCIAMENTO)\n")
	sb.WriteString("! ==============================================================================\n\n")

	for _, k := range ponKeys {
		ponONUs := ponMap[k]
		slot := k.slot
		port := k.port

		sb.WriteString(fmt.Sprintf("! ==============================================================================\n"))
		sb.WriteString(fmt.Sprintf("! PORTA PON %d/%d (%d ONUs)\n", slot, port, len(ponONUs)))
		sb.WriteString(fmt.Sprintf("! ==============================================================================\n\n"))

		// --- SEÇÃO 1: AUTORIZAÇÃO (Executada no modo cd onu) ---
		sb.WriteString(fmt.Sprintf("! --- 1. AUTORIZACAO DAS ONUS (PON %d/%d) ---\n", slot, port))
		sb.WriteString("cd onu\n")
		for _, onu := range ponONUs {
			serial := onu.SerialNumber
			if serial == "" {
				serial = "FHTT00000000"
			}
			clientDesc := onu.PPPoEUser
			if clientDesc == "" {
				clientDesc = "Bridge"
			}
			sb.WriteString(fmt.Sprintf("! ONU %d (%s)\n", onu.OnuID, clientDesc))
			sb.WriteString(fmt.Sprintf("set whitelist phy_addr address %s password null action add slot %d pon %d onu null type null\n",
				serial, slot, port))
		}
		sb.WriteString("\n")

		// --- SEÇÃO 2: WAN / INTERNET (Executada no modo cd lan) ---
		sb.WriteString(fmt.Sprintf("! --- 2. CONFIGURACAO DE WAN / INTERNET (PON %d/%d) ---\n", slot, port))
		sb.WriteString("cd lan\n\n")
		for _, onu := range ponONUs {
			vlan := 100 // default router
			if onu.PPPoEUser == "" {
				vlan = 10 // default bridge
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
				sb.WriteString(fmt.Sprintf("set wancfg sl %d %d %d ind 1 mode tr069_int ty r %d 65535 nat en qos dis vlanm tag tvlan dis 65535 65535 dsp pppoe pro dis %s %s null auto entries 6 fe1 fe2 fe3 fe4 ssid1 ssid5 \n",
					slot, port, onu.OnuID, vlan, login, senha))
				sb.WriteString(fmt.Sprintf("set wancfg sl %d %d %d ind 1 ip-stack-mode both ipv6-src-type dhcpv6 prefix-src-type delegate \n",
					slot, port, onu.OnuID))
				sb.WriteString(fmt.Sprintf("apply wancfg slot %d %d %d\n\n",
					slot, port, onu.OnuID))
			} else {
				sb.WriteString(fmt.Sprintf("! ONU Bridge %d - LAN\n", onu.OnuID))
				sb.WriteString(fmt.Sprintf("set epon slot %d pon %d onu %d port 1 service number 1\n",
					slot, port, onu.OnuID))
				sb.WriteString(fmt.Sprintf("set epon slot %d pon %d onu %d port 1 service 1 vlan_mode tag 0 33024 %d\n",
					slot, port, onu.OnuID, vlan))
				sb.WriteString(fmt.Sprintf("set epon slot %d pon %d onu %d port 1\n",
					slot, port, onu.OnuID))
				sb.WriteString(fmt.Sprintf("apply onu %d %d %d vlan\n\n",
					slot, port, onu.OnuID))
			}
		}

		// --- SEÇÃO 3: GERENCIAMENTO E CONTROLE (Volta com cd .. para o modo onu) ---
		sb.WriteString(fmt.Sprintf("! --- 3. GERENCIAMENTO E CONTROLE (PON %d/%d) ---\n", slot, port))
		sb.WriteString("cd ..\n")
		for _, onu := range ponONUs {
			if onu.PPPoEUser != "" {
				sb.WriteString(fmt.Sprintf("! ONU %d (%s) - Gerenciamento\n", onu.OnuID, onu.PPPoEUser))
				sb.WriteString(fmt.Sprintf("set onu_local_manage_con slot %d pon %d onu %d conf en cons en tel dis web en web_p 1025 web_ani_s en tel_ani_s dis web_admin_switch dis\n",
					slot, port, onu.OnuID))
				sb.WriteString(fmt.Sprintf("set remote_manage_cfg slot %d pon %d onu %d tr069 enable acs_url http://cwmp.nicnet.com.br:8088 acl_user Admin acl_pswd Admin@1234 inform enable interval 43200 port 30005 user cpe pswd cpe\n\n",
					slot, port, onu.OnuID))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
