package parser

import (
	"strings"
	"testing"
	"olt-migrate-backend/internal/models"
)

func TestParseFiberhome5000(t *testing.T) {
	config := `
!onu authorization configuration 
set autho sl 1 p 1 ty 5506-04-FA o 2 phy FHTT9a980b98 pas null 
set autho sl 1 p 2 ty 5506-01-A1 o 10 phy FHTT92370718 pas null

set wancfg sl 1 1 2 ind 1 mode tr069_in ty r 81 65535 nat en qos dis vlanm tag tvlan dis 65535 65535 dsp pppoe pro dis cliente.pppoe key:+906=mlk null auto entries 6 fe1
set wifi_serv_wlan slot 1 pon 1 onu 2 serv_no 1 index 1 ssid enable 2Alto do Castelo hide disable authmode wpa_psk/wpa2psk encrypt_type tkipaes wpakey @ltocastelo2 interval 0
`
	reader := strings.NewReader(config)
	parserInstance := &FiberhomeParser{}
	onus, err := parserInstance.Parse(reader)
	if err != nil {
		t.Fatalf("Erro ao parsear: %v", err)
	}

	if len(onus) != 2 {
		t.Fatalf("Esperava 2 ONUs, obteve %d", len(onus))
	}

	// Como a ordem não é garantida devido ao map, vamos iterar para achar a ONU 2
	var onu2 *models.ONU
	for _, o := range onus {
		if o.OnuID == 2 && o.SlotID == 1 && o.PortID == 1 {
			onu2 = &o
			break
		}
	}

	if onu2 == nil {
		t.Fatal("ONU 2 não encontrada")
	}

	if onu2.SerialNumber != "FHTT9a980b98" {
		t.Errorf("ONU 2 SN incorreto: %s", onu2.SerialNumber)
	}
	if onu2.PPPoEUser != "cliente.pppoe" || onu2.PPPoEPass != "senha123" {
		t.Errorf("PPPoE incorreto: %s / %s", onu2.PPPoEUser, onu2.PPPoEPass)
	}
	if onu2.WiFiSSID != "2Alto do Castelo" || onu2.WiFiPass != "@ltocastelo2" {
		t.Errorf("WiFi incorreto: %s / %s", onu2.WiFiSSID, onu2.WiFiPass)
	}
}

func TestParseFiberhome6000(t *testing.T) {
	config := `
authorize 1/1/1 1 type 5506-04-FA  phy-id FHTT9df6baa0 password null 
authorize 1/2/3 15 type HG6145E  phy-id FHTT9c0a4860 password null 

interface pon 1/1/1 
onu wan-cfg 1 ind 1 mode tr069-in ty r 189 0 nat en qos dis vlanm tag tvlan dis 0 0 dsp pppoe pro dis maria.silva key:mlkjih null pay upnp_switch
onu wifi connection 1 serv-no 1 index 1 ssid enable Maria Wi Fi hide disable authmode wpa-psk/wpa2psk encrypt-type tkipaes wpakey senha654 interval 0
`
	reader := strings.NewReader(config)
	parserInstance := &FiberhomeParser{}
	onus, err := parserInstance.Parse(reader)
	if err != nil {
		t.Fatalf("Erro ao parsear: %v", err)
	}

	if len(onus) != 2 {
		t.Fatalf("Esperava 2 ONUs, obteve %d", len(onus))
	}

	var onu1 *models.ONU
	for _, o := range onus {
		if o.OnuID == 1 && o.SlotID == 1 && o.PortID == 1 {
			onu1 = &o
			break
		}
	}

	if onu1 == nil {
		t.Fatal("ONU 1 não encontrada")
	}

	if onu1.SerialNumber != "FHTT9df6baa0" {
		t.Errorf("ONU 1 SN incorreto: %s", onu1.SerialNumber)
	}
	if onu1.PPPoEUser != "maria.silva" || onu1.PPPoEPass != "123456" {
		t.Errorf("PPPoE incorreto: %s / %s", onu1.PPPoEUser, onu1.PPPoEPass)
	}
	if onu1.WiFiSSID != "Maria Wi Fi" || onu1.WiFiPass != "senha654" {
		t.Errorf("WiFi incorreto: %s / %s", onu1.WiFiSSID, onu1.WiFiPass)
	}
}

func TestParseFiberhome_EmilioCase(t *testing.T) {
	config := `
set autho sl 1 p 1 ty 5506-04-FA o 6 phy FHTT9a980b98 pas null 
set wancfg sl 1 1 6 ind 1 mode tr069_in ty r 81 65535 nat en qos dis vlanm tag tvlan dis 65535 65535 dsp pppoe pro dis 118985.emilio.0 key:+906=mlk null auto entries 6 fe1
set wifi_serv_wlan slot 1 pon 1 onu 6 serv_no 1 index 1 ssid enable Redenicnet hide disable authmode wpa_psk/wpa2psk encrypt_type tkipaes wpakey saocarlos interval 0
set wifi_serv_wlan slot 1 pon 1 onu 6 serv_no 2 index 1 ssid enable Redenicnet_5G hide disable authmode wpa_psk/wpa2psk encrypt_type tkipaes wpakey saocarlos interval 0
`
	reader := strings.NewReader(config)
	parserInstance := &FiberhomeParser{}
	onus, err := parserInstance.Parse(reader)
	if err != nil {
		t.Fatalf("Erro ao parsear: %v", err)
	}
	if len(onus) != 1 {
		t.Fatalf("Esperava 1 ONU, obteve %d", len(onus))
	}
	if onus[0].PPPoEUser != "118985.emilio.0" {
		t.Errorf("PPPoEUser incorreto: %s", onus[0].PPPoEUser)
	}
	if onus[0].WiFiSSID != "Redenicnet" || onus[0].WiFiPass != "saocarlos" {
		t.Errorf("WiFi incorreto, esperado Redenicnet/saocarlos, obteve: %s / %s", onus[0].WiFiSSID, onus[0].WiFiPass)
	}
}

func TestParseFiberhome_PositionalAndOutOfOrderWiFi(t *testing.T) {
	// Teste com sintaxe posicional sl 1 1 6 e wifi vindo antes de autho (ou sem palavras-chave explícitas)
	config := `
set wifi_serv_wlan sl 1 1 6 serv_no 1 index 1 ssid enable Positional_WiFi hide disable authmode wpa_psk/wpa2psk encrypt_type tkipaes wpakey segredo123 interval 0
set autho sl 1 p 1 ty 5506-04-FA o 6 phy FHTT9a980b98 pas null
set wancfg sl 1 1 6 ind 1 mode tr069_in ty r 81 65535 nat en qos dis vlanm tag tvlan dis 65535 65535 dsp pppoe pro dis positional.user key:+906=mlk null auto entries 6 fe1
`
	reader := strings.NewReader(config)
	parserInstance := &FiberhomeParser{}
	onus, err := parserInstance.Parse(reader)
	if err != nil {
		t.Fatalf("Erro ao parsear: %v", err)
	}
	if len(onus) != 1 {
		t.Fatalf("Esperava 1 ONU, obteve %d", len(onus))
	}
	if onus[0].WiFiSSID != "Positional_WiFi" || onus[0].WiFiPass != "segredo123" {
		t.Errorf("WiFi incorreto, esperado Positional_WiFi/segredo123, obteve: %s / %s", onus[0].WiFiSSID, onus[0].WiFiPass)
	}
}

func TestParseFiberhome_UserCase(t *testing.T) {
	config := `
set wancfg sl 1 8 29 ind 1 mode tr069_in ty r 81 65535 nat en qos dis vlanm tag tvlan dis 65535 65535 dsp pppoe pro dis 93620.linx.2 key:kkkeelm; null auto entries 6 fe1 fe2 fe3 fe4 ssid1 ssid5 

set wifi_serv_wlan slot 1 pon 8 onu 29 serv_no 1 index 1 ssid enable Ellos hide disable authmode wpa_psk/wpa2psk encrypt_type tkipaes wpakey 20112024 interval 86400 radius_serv ipv4 0.0.0.0 port 0 pswd null wep_length 40bit key_index 1 wep_key 12345 12345 12345 12345 wapi_serv_addr 0.0.0.0 0 wifi_connect_num 32 

set wifi_serv_cfg slot 1 pon 8 onu 29 serv_no 1 wifi disable district etsi channel 0 standard 802.11b txpower 0 frequency 2.4ghz freq_bandwidth 20mhz/40mhz
set wifi_serv_wlan slot 1 pon 8 onu 29 serv_no 2 index 1 ssid enable Ellos_5G hide disable authmode wpa_psk/wpa2psk encrypt_type tkipaes wpakey 20112024 interval 86400 radius_serv ipv4 0.0.0.0 port 0 pswd null wep_length 40bit key_index 1 wep_key 12345 12345 12345 12345 wapi_serv_addr 0.0.0.0 0 wifi_connect_num 32
`
	reader := strings.NewReader(config)
	parserInstance := &FiberhomeParser{}
	onus, err := parserInstance.Parse(reader)
	if err != nil {
		t.Fatalf("Erro ao parsear: %v", err)
	}
	if len(onus) != 1 {
		t.Fatalf("Esperava 1 ONU, obteve %d", len(onus))
	}
	t.Logf("ONU parseada: %+v", onus[0])
	if onus[0].SlotID != 1 || onus[0].PortID != 8 || onus[0].OnuID != 29 {
		t.Errorf("Identificadores incorretos: slot=%d, port=%d, onuID=%d", onus[0].SlotID, onus[0].PortID, onus[0].OnuID)
	}
	if onus[0].WiFiSSID != "Ellos" || onus[0].WiFiPass != "20112024" {
		t.Errorf("WiFi incorreto, esperado Ellos/20112024, obteve: '%s' / '%s'", onus[0].WiFiSSID, onus[0].WiFiPass)
	}
}
