package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"olt-migrate-backend/internal/crypto"
	"olt-migrate-backend/internal/models"
)

// Regex para Fiberhome 5000
var re5000ONU = regexp.MustCompile(`(?:set\s+)?autho\s+(?:sl|slot)?\s*(\d+)\s+(?:p|pon)?\s*(\d+)\s+(?:ty|type)?\s*(\S+)\s+(?:o|onu)?\s*(\d+)\s+(?:phy|phy_addr)?\s*(\S+)`)
var re5000WanCfg = regexp.MustCompile(`(?:set\s+)?wancfg\s+(?:sl|slot)?\s*(\d+)\s+(?:p|pon)?\s*(\d+)\s+(?:o|onu)?\s*(\d+).*?dsp\s+pppoe.*?(?:pro\s+(?:dis|disable|en|enable)|dsp\s+pppoe)\s+(\S+)\s+(?:key:)?(\S+)`)
var re5000WiFiBase = regexp.MustCompile(`(?:set\s+)?(?:wifi|wlan)\S*\s+(?:slot|sl)?\s*(\d+)\s+(?:pon|p)?\s*(\d+)\s+(?:onu|o)?\s*(\d+)`)

// Regex para Fiberhome 6000
var re6000ONU = regexp.MustCompile(`authorize\s+\d+/(\d+)/(\d+)\s+(\d+)\s+type\s+(\S+)\s+phy-id\s+(\S+)`)
var re6000Interface = regexp.MustCompile(`^interface\s+pon\s+\d+/(\d+)/(\d+)`)
var re6000WanCfg = regexp.MustCompile(`^onu\s+wan-cfg\s+(\d+).*?dsp\s+pppoe.*?(?:pro\s+(?:dis|disable|en|enable)|dsp\s+pppoe)\s+(\S+)\s+(?:key:)?(\S+)`)
var re6000WiFiBase = regexp.MustCompile(`(?:^|\s)(?:onu\s+)?(?:wifi|wlan)\S*(?:\s+\S+)?\s+(\d+)`)

func makeKey(slot, port, onuID int) string {
	return fmt.Sprintf("%d-%d-%d", slot, port, onuID)
}

// Funções auxiliares flexíveis para capturar SSID e Chave WPA em qualquer linha WiFi, independente de ordem ou sufixos
func extractWiFiParams(line string) (string, string) {
	var ssid, pass string

	// 1. Captura SSID após "ssid enable ", "ssid en ", "ssid " de forma agnóstica à palavra-chave seguinte
	reSSID := regexp.MustCompile(`(?:ssid|ssid[-_]?name|ssid[-_]?cfg|ssid[-_]?val)\s+(?:enable|en|disable|dis|name)?\s*["']?([^\s"']+(?:\s+[^\s"']+)*?)["']?(?:\s+(?:hide|auth|wpa|enc|key|int|rad|wep|wap|wifi|serv|chan|stan|tx|freq|port|pswd|mode|vlan|eth|tag|pri|tpid|vid|dis|en|max|num|connect|val|idx|index|wmm|isolation|assoc|bssid|beacon|dtim|rts|frag)|$)`)
	if matches := reSSID.FindStringSubmatch(line); matches != nil {
		ssid = strings.TrimSpace(matches[1])
	} else {
		// Fallback ultra-simples para pegar qualquer coisa entre ssid [enable] e a próxima palavra ou aspas
		reSSIDFallback := regexp.MustCompile(`(?:ssid|ssid[-_]?name)\s+(?:enable|en|disable|dis)?\s*["']?([^"'\s]+)["']?`)
		if m := reSSIDFallback.FindStringSubmatch(line); m != nil {
			ssid = strings.TrimSpace(m[1])
		}
	}

	if ssid != "" {
		// Remove sufixos de rádio para manter o nome base puro e unificado
		ssid = strings.TrimSuffix(ssid, "_5G")
		ssid = strings.TrimSuffix(ssid, "_5g")
		ssid = strings.TrimSuffix(ssid, "-5G")
		ssid = strings.TrimSuffix(ssid, "_2.4G")
		ssid = strings.TrimSuffix(ssid, "-2.4G")
	}

	// 2. Captura WPA Key em wpakey, wpa-key, wpa_key, pswd, password, wpa_pswd, wpa-pswd, wpa_pass
	reKey := regexp.MustCompile(`(?:wpakey|wpa[-_]?key|pswd|password|wpa[-_]?pswd|wpa[-_]?pass|key[-_]?val|wpa[-_]?key[-_]?val|key)\s+(?:key:)?(["']?\S+?["']?)(?:\s|$)`)
	if matches := reKey.FindStringSubmatch(line); matches != nil {
		pass = strings.TrimSpace(matches[1])
		pass = strings.Trim(pass, `"`)
		pass = strings.Trim(pass, `'`)
		if strings.HasPrefix(pass, "key:") || strings.Contains(line, "key:"+pass) {
			pass = crypto.DecryptPassword(pass)
		}
	}

	return ssid, pass
}

// extractSlotPortOnuID identifica (slot, port, onuID) de qualquer linha Fiberhome independente de ordem ou formato
func extractSlotPortOnuID(line string, currentSlot, currentPort int) (int, int, int, bool) {
	// 1. Formato explícito: slot/sl X pon/p Y onu/onuid/o/id Z (em qualquer lugar da linha, mesmo com palavras no meio)
	reExplicit := regexp.MustCompile(`(?:slot|sl)\s+(\d+).*?(?:pon|p)\s+(\d+).*?(?:onu|onuid|o|id)\s+(\d+)`)
	if m := reExplicit.FindStringSubmatch(line); m != nil {
		s, _ := strconv.Atoi(m[1])
		p, _ := strconv.Atoi(m[2])
		o, _ := strconv.Atoi(m[3])
		return s, p, o, true
	}

	// 2. Formato posicional AN5000: sl/slot X Y Z (ex: sl 1 1 6 ou slot 1 1 6)
	rePositional := regexp.MustCompile(`(?:^|\s)(?:slot|sl)\s+(\d+)\s+(\d+)\s+(\d+)(?:\s|$)`)
	if m := rePositional.FindStringSubmatch(line); m != nil {
		s, _ := strconv.Atoi(m[1])
		p, _ := strconv.Atoi(m[2])
		o, _ := strconv.Atoi(m[3])
		return s, p, o, true
	}

	// 3. Formato com barra AN6000: 1/1/1 6 ou 0/1/1 6
	reSlash := regexp.MustCompile(`(?:^|\s)(?:\d+/)?(\d+)/(\d+)\s+(?:onu|onuid|o|id)?\s*(\d+)(?:\s|$)`)
	if m := reSlash.FindStringSubmatch(line); m != nil {
		s, _ := strconv.Atoi(m[1])
		p, _ := strconv.Atoi(m[2])
		o, _ := strconv.Atoi(m[3])
		return s, p, o, true
	}

	// 4. Formato contextual AN6000 sob interface pon X/Y: apenas o número do ONU ID após o comando
	if currentSlot > 0 && currentPort > 0 {
		reSingleOnu := regexp.MustCompile(`(?:^|\s)(?:connection|attribute|wlan|wifi|cfg|wan-cfg|wancfg|onu|wan|service|profile|port|port-isolation|remote-manage-cfg|authorize|whitelist|vlan)\s+(\d+)(?:\s|$)`)
		if m := reSingleOnu.FindStringSubmatch(line); m != nil {
			o, _ := strconv.Atoi(m[1])
			return currentSlot, currentPort, o, true
		}
	}

	return 0, 0, 0, false
}

func getOrCreateONU(onuMap map[string]*models.ONU, slot, port, onuID int) *models.ONU {
	key := makeKey(slot, port, onuID)
	onu, exists := onuMap[key]
	if !exists {
		onu = &models.ONU{
			SlotID: slot,
			PortID: port,
			OnuID:  onuID,
		}
		onuMap[key] = onu
	}
	return onu
}

type FiberhomeParser struct{}

// Parse lê o arquivo e retorna uma lista de ONUs configuradas com resiliência total a ordens e sintaxes
func (p *FiberhomeParser) Parse(reader io.Reader) ([]models.ONU, error) {
	onuMap := make(map[string]*models.ONU)
	scanner := bufio.NewScanner(reader)
	const maxCapacity = 10 * 1024 * 1024 // 10 MB per line
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, maxCapacity)

	currentSlot := 0
	currentPort := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "!") || line == "" {
			continue
		}

		// Rastreia contexto de interface da série 6000 (ex: interface pon 1/1/8)
		if matches := re6000Interface.FindStringSubmatch(line); matches != nil {
			currentSlot, _ = strconv.Atoi(matches[1])
			currentPort, _ = strconv.Atoi(matches[2])
			continue
		}
		if line == "exit" || line == "cd .." || strings.HasPrefix(line, "interface ") || strings.HasPrefix(line, "cd ") {
			if strings.HasPrefix(line, "interface ") && !strings.Contains(line, "pon") {
				currentSlot = 0
				currentPort = 0
			}
		}

		// Tenta identificar (slot, port, onuID) na linha via motor universal
		slot, port, onuID, ok := extractSlotPortOnuID(line, currentSlot, currentPort)
		if !ok {
			// Fallback com regex legados
			if matches := re5000ONU.FindStringSubmatch(line); matches != nil {
				slot, _ = strconv.Atoi(matches[1])
				port, _ = strconv.Atoi(matches[2])
				onuID, _ = strconv.Atoi(matches[4])
				ok = true
			} else if matches := re6000ONU.FindStringSubmatch(line); matches != nil {
				slot, _ = strconv.Atoi(matches[1])
				port, _ = strconv.Atoi(matches[2])
				onuID, _ = strconv.Atoi(matches[3])
				ok = true
			}
		}

		if ok && slot > 0 && port > 0 && onuID > 0 {
			onu := getOrCreateONU(onuMap, slot, port, onuID)

			// 1. Extração de Serial Number e Modelo (Autorização)
			if strings.Contains(line, "autho") || strings.Contains(line, "authorize") || strings.Contains(line, "whitelist") || strings.Contains(line, "phy") || strings.Contains(line, "sn ") || strings.Contains(line, "mac ") {
				reSN := regexp.MustCompile(`(?:phy|phy[-_]?id|phy[-_]?addr|sn|serial|mac)\s+([A-Za-z0-9-_]{8,16})`)
				if m := reSN.FindStringSubmatch(line); m != nil {
					sn := strings.TrimSpace(m[1])
					if sn != "null" && sn != "00000000" && sn != "FHTT00000000" {
						onu.SerialNumber = sn
					}
				} else if onu.SerialNumber == "" {
					reSNFallback := regexp.MustCompile(`\b([A-Z]{4}[0-9a-fA-F]{8})\b`)
					if m := reSNFallback.FindStringSubmatch(line); m != nil {
						if m[1] != "FHTT00000000" {
							onu.SerialNumber = m[1]
						}
					}
				}

				reModel := regexp.MustCompile(`(?:type|ty|model)\s+([A-Za-z0-9_-]+)`)
				if m := reModel.FindStringSubmatch(line); m != nil {
					model := strings.TrimSpace(m[1])
					if model != "null" && model != "" {
						onu.Model = model
					}
				}
			}

			// 2. Extração de WAN / PPPoE Credentials
			if strings.Contains(line, "wancfg") || strings.Contains(line, "wan-cfg") || strings.Contains(line, "wan ") || strings.Contains(line, "pppoe") || strings.Contains(line, "key:") {
				rePass := regexp.MustCompile(`(?:pswd|pass|password|key)\s+(?:key:)?([^\s"]+)`)
				if m := rePass.FindStringSubmatch(line); m != nil {
					pass := strings.TrimSpace(m[1])
					if pass != "null" && pass != "auto" {
						if strings.HasPrefix(m[0], "key:") || strings.Contains(line, "key:"+pass) || strings.Contains(pass, "=") || strings.Contains(pass, "<") || strings.Contains(pass, ";") {
							pass = crypto.DecryptPassword(pass)
						}
						onu.PPPoEPass = pass
					}
				} else {
					reKeyDirect := regexp.MustCompile(`key:([^\s"]+)`)
					if m := reKeyDirect.FindStringSubmatch(line); m != nil {
						onu.PPPoEPass = crypto.DecryptPassword(m[1])
					}
				}

				reUser := regexp.MustCompile(`(?:pro\s+(?:dis|disable|en|enable)|user|username|pppoe-user)\s+([^\s:"]+)`)
				if m := reUser.FindStringSubmatch(line); m != nil {
					user := strings.TrimSpace(m[1])
					if user != "null" && user != "auto" && user != "dis" && user != "en" && user != "dis," && user != "en," {
						onu.PPPoEUser = user
					}
				} else if strings.Contains(line, "key:") {
					parts := strings.Fields(line)
					for idx, p := range parts {
						if strings.HasPrefix(p, "key:") && idx > 0 {
							candidate := parts[idx-1]
							if candidate != "dis" && candidate != "enable" && candidate != "en" && candidate != "null" && candidate != "auto" {
								onu.PPPoEUser = candidate
								break
							}
						}
					}
				} else if m := re5000WanCfg.FindStringSubmatch(line); m != nil {
					onu.PPPoEUser = m[4]
					if onu.PPPoEPass == "" {
						onu.PPPoEPass = crypto.DecryptPassword(m[5])
					}
				} else if m := re6000WanCfg.FindStringSubmatch(line); m != nil {
					onu.PPPoEUser = m[2]
					if onu.PPPoEPass == "" {
						onu.PPPoEPass = crypto.DecryptPassword(m[3])
					}
				}
			}

			// 3. Extração de Wi-Fi
			if strings.Contains(line, "wifi") || strings.Contains(line, "wlan") || strings.Contains(line, "ssid") || strings.Contains(line, "wpakey") || strings.Contains(line, "pswd") {
				ssid, pass := extractWiFiParams(line)
				is5G := strings.Contains(line, "serv_no 2") || strings.Contains(line, "serv-no 2") || strings.Contains(line, "serv 2") || strings.Contains(line, "ssid 2") || strings.Contains(line, "index 2") ||
					strings.Contains(line, "serv_no 5") || strings.Contains(line, "serv-no 5") || strings.Contains(line, "serv 5") || strings.Contains(line, "ssid 5") || strings.Contains(line, "index 5") ||
					strings.Contains(line, "_5G") || strings.Contains(line, "-5G") || strings.Contains(line, "5.8ghz")

				if ssid != "" {
					if !is5G || onu.WiFiSSID == "" {
						onu.WiFiSSID = ssid
					}
				}
				if pass != "" {
					if !is5G || onu.WiFiPass == "" {
						onu.WiFiPass = pass
					}
				}
			}

			// 4. Extração de VLAN
			if strings.Contains(line, "vlan") || strings.Contains(line, "vid ") || strings.Contains(line, "cvlan") || strings.Contains(line, "tpid") {
				reVLAN := regexp.MustCompile(`(?:vlan|vid|tag|cvlan|vlanm)\s+(?:tag\s+|vid\s+|priority\s+\d+\s+tpid\s+\d+\s+vid\s+)?(\d+)`)
				if m := reVLAN.FindStringSubmatch(line); m != nil {
					vlan, _ := strconv.Atoi(m[1])
					if vlan > 0 && vlan <= 4094 {
						if len(onu.Services) == 0 {
							onu.Services = append(onu.Services, models.Service{CVLAN: vlan})
						} else {
							onu.Services[0].CVLAN = vlan
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var onus []models.ONU
	for _, onu := range onuMap {
		onus = append(onus, *onu)
	}

	return onus, nil
}
