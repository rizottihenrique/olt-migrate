package generator

import (
	"fmt"
	"sort"
	"strings"

	"olt-migrate-backend/internal/models"
)

// Fiberhome6000Generator gera comandos hierárquicos para OLTs linha AN6000, agrupados por PON
type Fiberhome6000Generator struct{}

// Generate recebe as ONUs mapeadas e retorna o script CLI no formato AN6000, organizado por interface PON
func (g *Fiberhome6000Generator) Generate(onus []models.ONU) string {
	if len(onus) == 0 {
		return "! Nenhuma ONU para migrar.\n"
	}

	// 1. Ordenar as ONUs de forma determinística por Slot, Port e OnuID
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

	var sb strings.Builder

	// --- SEÇÃO 1: COMANDOS GLOBAIS DE AUTORIZAÇÃO ---
	sb.WriteString("! ==================================================\n")
	sb.WriteString("! 1. AUTORIZACAO GLOBAIS DE ONUS\n")
	sb.WriteString("! ==================================================\n\n")

	for _, onu := range sortedONUs {
		model := onu.Model
		if model == "" {
			model = "5506-04-FA"
		}
		serial := onu.SerialNumber
		if serial == "" {
			serial = "FHTT00000000"
		}

		if onu.PPPoEUser != "" {
			// Router
			sb.WriteString(fmt.Sprintf("whitelist add phy-id %s checkcode fiberhome type %s slot %d pon %d onuid %d\n",
				serial, model, onu.SlotID, onu.PortID, onu.OnuID))
			sb.WriteString(fmt.Sprintf("authorize 1/%d/%d %d type %s phy-id %s password null\n",
				onu.SlotID, onu.PortID, onu.OnuID, model, serial))
			sb.WriteString(fmt.Sprintf("onu pon-type 1/%d/%d %d 712\n\n",
				onu.SlotID, onu.PortID, onu.OnuID))
		} else {
			// Bridge
			sb.WriteString(fmt.Sprintf("whitelist add phy-id %s type null slot %d pon %d onuid %d\n",
				serial, onu.SlotID, onu.PortID, onu.OnuID))
			sb.WriteString(fmt.Sprintf("authorize 1/%d/%d %d type null phy-id %s password null\n",
				onu.SlotID, onu.PortID, onu.OnuID, serial))
			sb.WriteString(fmt.Sprintf("onu pon-type 1/%d/%d %d 712\n\n",
				onu.SlotID, onu.PortID, onu.OnuID))
		}
	}

	// --- SEÇÃO 2: PROVISIONAMENTO AGRUPADO POR PON ---
	sb.WriteString("! ==================================================\n")
	sb.WriteString("! 2. MIGRACAO PARA FIBERHOME AN6000 (SEPARADO POR PON)\n")
	sb.WriteString("! ==================================================\n\n")

	currentSlot := -1
	currentPort := -1

	for _, onu := range sortedONUs {
		vlan := 188 // default router vlan
		if onu.PPPoEUser == "" {
			vlan = 10 // default bridge vlan
		}
		if len(onu.Services) > 0 && onu.Services[0].CVLAN > 0 {
			vlan = onu.Services[0].CVLAN
		}

		// Mudança de interface PON: fecha a anterior e abre a nova
		if onu.SlotID != currentSlot || onu.PortID != currentPort {
			if currentSlot != -1 {
				sb.WriteString("exit\n\n")
			}
			currentSlot = onu.SlotID
			currentPort = onu.PortID
			sb.WriteString(fmt.Sprintf("interface pon 1/%d/%d\n\n", currentSlot, currentPort))
		}

		if onu.PPPoEUser != "" {
			// --- ONT ROUTER AN6000 ---
			login := onu.PPPoEUser
			senha := onu.PPPoEPass
			if senha == "" {
				senha = "123456"
			}
			ssid := onu.WiFiSSID
			if ssid == "" {
				ssid = login
			}
			wifiPass := onu.WiFiPass
			if wifiPass == "" {
				wifiPass = "12345678"
			}

			sb.WriteString(fmt.Sprintf("! --- ONU %d (%s) ---\n", onu.OnuID, login))
			sb.WriteString(fmt.Sprintf("onu wan-cfg %d ind 1 mode tr069-int ty r %d 65535 nat en qos dis vlan tag tvlan dis 65535 65535 dsp pppoe pro dis %s %s null auto upnp_switch enable entries 6 fe1 fe2 fe3 fe4 ssid1 ssid5\n",
				onu.OnuID, vlan, login, senha))
			sb.WriteString(fmt.Sprintf("onu ipv6-wan-cfg %d ind 1 ip-stack-mode both ipv6-src-type dhcpv6 prefix-src-type delegate ipv6-address ::/0 ipv6-gateway :: ipv6-master-dns :: ipv6-slave-dns :: ipv6-static-prefix ::/0\n",
				onu.OnuID))
			sb.WriteString(fmt.Sprintf("onu wifi connection %d serv-no 1 index 1 ssid enable %s hide disable authmode wpa-psk/wpa2psk encrypt-type tkipaes wpakey %s interval 0 radius-serv ipv4 192.168.1.18 port 1812 pswd 12345678 wep-length 40bit key-index 1 wep-key 12345 12345 12345 12345 wapi-serv-addr 0.0.0.0 0 wifi-connect-num 32\n",
				onu.OnuID, ssid, wifiPass))
			sb.WriteString(fmt.Sprintf("onu wifi attribute %d serv-no 1 wifi enable district brazil channel 0 standard 802.11bgn txpower 20 frequency 2.4ghz freq-bandwidth 20mhz\n",
				onu.OnuID))
			sb.WriteString(fmt.Sprintf("onu wifi connection %d serv-no 2 index 1 ssid enable %s_5G hide disable authmode wpa-psk/wpa2psk encrypt-type tkipaes wpakey %s interval 0 radius-serv ipv4 192.168.1.18 port 1812 pswd 12345678 wep-length 40bit key-index 1 wep-key 12345 12345 12345 12345 wapi-serv-addr 0.0.0.0 0 wifi-connect-num 32\n",
				onu.OnuID, ssid, wifiPass))
			sb.WriteString(fmt.Sprintf("onu wifi attribute %d serv-no 2 wifi enable district brazil channel 0 standard 802.11ac txpower 20 frequency 5.8ghz freq-bandwidth 80mhz\n",
				onu.OnuID))
			sb.WriteString(fmt.Sprintf("onu port-isolation disable %d\n",
				onu.OnuID))
			sb.WriteString(fmt.Sprintf("onu remote-manage-cfg %d tr069 enable acs-url http://cwmp.nicnet.com.br:8088 acl-user Admin acl-pswd Admin@1234 inform enable interval 59834 port 0 user Admin pswd Admin@1234\n\n",
				onu.OnuID))
		} else {
			// --- ONU BRIDGE AN6000 ---
			sb.WriteString(fmt.Sprintf("! --- ONU Bridge %d ---\n", onu.OnuID))
			sb.WriteString(fmt.Sprintf("onu port vlan %d eth 1 service count 1\n",
				onu.OnuID))
			sb.WriteString(fmt.Sprintf("onu port vlan %d eth 1 service 1 tag priority 255 tpid 33024 vid %d\n\n",
				onu.OnuID, vlan))
		}
	}

	if currentSlot != -1 {
		sb.WriteString("exit\n")
	}

	return sb.String()
}
