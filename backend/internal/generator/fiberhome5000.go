package generator

import (
	"fmt"
	"sort"
	"strings"

	"olt-migrate-backend/internal/models"
)

// Fiberhome5000Generator gera comandos para OLTs linha AN5000, organizados por PON
type Fiberhome5000Generator struct{}

// Generate recebe as ONUs mapeadas e retorna o script CLI no formato AN5000
func (g *Fiberhome5000Generator) Generate(onus []models.ONU) string {
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

	sb.WriteString("! ==================================================\n")
	sb.WriteString("! MIGRACAO PARA FIBERHOME AN5000 (SEPARADO POR PON)\n")
	sb.WriteString("! ==================================================\n\n")

	currentSlot := -1
	currentPort := -1

	for _, onu := range sortedONUs {
		serial := onu.SerialNumber
		if serial == "" {
			serial = "FHTT00000000"
		}

		vlan := 100 // default router
		if onu.PPPoEUser == "" {
			vlan = 10 // default bridge
		}
		if len(onu.Services) > 0 && onu.Services[0].CVLAN > 0 {
			vlan = onu.Services[0].CVLAN
		}

		if onu.SlotID != currentSlot || onu.PortID != currentPort {
			currentSlot = onu.SlotID
			currentPort = onu.PortID
			sb.WriteString(fmt.Sprintf("! --- PORTA PON %d/%d ---\n\n", currentSlot, currentPort))
		}

		if onu.PPPoEUser != "" {
			// --- ONT ROUTER AN5000 ---
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

			sb.WriteString(fmt.Sprintf("! ONU %d (%s)\n", onu.OnuID, login))
			sb.WriteString("cd onu\n")
			sb.WriteString(fmt.Sprintf("set whitelist phy_addr address %s password null action add slot %d pon %d onu null type null\n\n",
				serial, onu.SlotID, onu.PortID))
			sb.WriteString("cd onu\n")
			sb.WriteString("cd lan\n\n")
			sb.WriteString(fmt.Sprintf("set wancfg sl %d %d %d ind 1 mode tr069_int ty r %d 65535 nat en qos dis vlanm tag tvlan dis 65535 65535 dsp pppoe pro dis %s %s null auto entries 6 fe1 fe2 fe3 fe4 ssid1 ssid5 \n\n",
				onu.SlotID, onu.PortID, onu.OnuID, vlan, login, senha))
			sb.WriteString(fmt.Sprintf("set wancfg sl %d %d %d ind 1 ip-stack-mode both ipv6-src-type dhcpv6 prefix-src-type delegate \n\n",
				onu.SlotID, onu.PortID, onu.OnuID))
			sb.WriteString(fmt.Sprintf("apply wancfg slot %d %d %d\n\n",
				onu.SlotID, onu.PortID, onu.OnuID))
			sb.WriteString(fmt.Sprintf("set wifi_serv_wlan slot %d pon %d onu %d serv_no 1 index 1 ssid enable %s hide disable authmode wpa_psk/wpa2psk encrypt_type tkipaes wpakey %s interval 0 radius_serv ipv4 192.168.1.18 port 1812 pswd 12345678 wep_length 40bit key_index 1 wep_key 12345 12345 12345 12345 wapi_serv_addr 0.0.0.0 0 wifi_connect_num 32\n\n",
				onu.SlotID, onu.PortID, onu.OnuID, ssid, wifiPass))
			sb.WriteString(fmt.Sprintf("set wifi_serv_cfg slot %d pon %d onu %d serv_no 1 wifi enable district brazil channel 0 standard 802.11bgn txpower 20 frequency 2.4ghz freq_bandwidth 20mhz\n\n",
				onu.SlotID, onu.PortID, onu.OnuID))
			sb.WriteString(fmt.Sprintf("set wifi_serv_wlan slot %d pon %d onu %d serv_no 2 index 1 ssid enable %s_5G hide disable authmode wpa_psk/wpa2psk encrypt_type tkipaes wpakey %s interval 0 radius_serv ipv4 192.168.1.18 port 1812 pswd 12345678 wep_length 40bit key_index 1 wep_key 12345 12345 12345 12345 wapi_serv_addr 0.0.0.0 0 wifi_connect_num 32\n\n",
				onu.SlotID, onu.PortID, onu.OnuID, ssid, wifiPass))
			sb.WriteString(fmt.Sprintf("set wifi_serv_cfg slot %d pon %d onu %d serv_no 2 wifi enable district brazil channel 0 standard 802.11ac txpower 20 frequency 5.8ghz freq_bandwidth 80mhz\n\n",
				onu.SlotID, onu.PortID, onu.OnuID))
			sb.WriteString("cd ..\n\n")
			sb.WriteString(fmt.Sprintf("set onu_local_manage_con slot %d pon %d onu %d conf en cons en tel dis web en web_p 1025 web_ani_s en tel_ani_s dis web_admin_switch dis\n\n",
				onu.SlotID, onu.PortID, onu.OnuID))
			sb.WriteString(fmt.Sprintf("set remote_manage_cfg slot %d pon %d onu %d tr069 enable acs_url http://cwmp.nicnet.com.br:8088 acl_user Admin acl_pswd Admin@1234 inform enable interval 43200 port 30005 user cpe pswd cpe\n\n",
				onu.SlotID, onu.PortID, onu.OnuID))
		} else {
			// --- ONU BRIDGE AN5000 ---
			sb.WriteString(fmt.Sprintf("! ONU Bridge %d\n", onu.OnuID))
			sb.WriteString("cd onu\n")
			sb.WriteString(fmt.Sprintf("set whitelist phy_addr address %s password null action add slot %d pon %d onu null type null\n\n",
				serial, onu.SlotID, onu.PortID))
			sb.WriteString("cd onu\n")
			sb.WriteString("cd lan\n")
			sb.WriteString(fmt.Sprintf("set epon slot %d pon %d onu %d port 1 service number 1\n",
				onu.SlotID, onu.PortID, onu.OnuID))
			sb.WriteString(fmt.Sprintf("set epon slot %d pon %d onu %d port 1 service 1 vlan_mode tag 0 33024 %d\n",
				onu.SlotID, onu.PortID, onu.OnuID, vlan))
			sb.WriteString(fmt.Sprintf("set epon slot %d pon %d onu %d port 1\n",
				onu.SlotID, onu.PortID, onu.OnuID))
			sb.WriteString(fmt.Sprintf("apply onu %d %d %d vlan\n\n",
				onu.SlotID, onu.PortID, onu.OnuID))
		}
	}

	return sb.String()
}
