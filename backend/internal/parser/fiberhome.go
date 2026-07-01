package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"olt-migrate-backend/internal/models"
)

// Regex para Fiberhome 5000
var re5000ONU = regexp.MustCompile(`set autho sl (\d+) p (\d+) ty (\S+) o (\d+) phy (\S+)`)
var re5000WanCfg = regexp.MustCompile(`set wancfg sl (\d+) (\d+) (\d+).*?dsp pppoe.*?dis (\S+)\s+key:(\S+)`)

// Regex para Fiberhome 6000
var re6000ONU = regexp.MustCompile(`authorize \d+/(\d+)/(\d+) (\d+) type (\S+)\s+phy-id (\S+)`)
var re6000Interface = regexp.MustCompile(`^interface pon \d+/(\d+)/(\d+)`)
var re6000WanCfg = regexp.MustCompile(`^onu wan-cfg (\d+).*?dsp pppoe.*?dis (\S+)\s+key:(\S+)`)

func makeKey(slot, port, onuID int) string {
	return fmt.Sprintf("%d-%d-%d", slot, port, onuID)
}

type FiberhomeParser struct{}

// Parse lê o arquivo e retorna uma lista de ONUs configuradas
func (p *FiberhomeParser) Parse(reader io.Reader) ([]models.ONU, error) {
	onuMap := make(map[string]*models.ONU)
	scanner := bufio.NewScanner(reader)

	currentSlot := 0
	currentPort := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "!") {
			continue
		}

		// --- Parsing 5000 Series ---
		if matches := re5000ONU.FindStringSubmatch(line); matches != nil {
			slot, _ := strconv.Atoi(matches[1])
			port, _ := strconv.Atoi(matches[2])
			model := matches[3]
			onuID, _ := strconv.Atoi(matches[4])
			sn := matches[5]

			key := makeKey(slot, port, onuID)
			onuMap[key] = &models.ONU{
				SlotID:       slot,
				PortID:       port,
				OnuID:        onuID,
				SerialNumber: sn,
				Model:        model,
			}
			continue
		}

		if matches := re5000WanCfg.FindStringSubmatch(line); matches != nil {
			slot, _ := strconv.Atoi(matches[1])
			port, _ := strconv.Atoi(matches[2])
			onuID, _ := strconv.Atoi(matches[3])
			user := matches[4]
			pass := matches[5]

			key := makeKey(slot, port, onuID)
			if onu, exists := onuMap[key]; exists {
				onu.PPPoEUser = user
				onu.PPPoEPass = pass
			}
			continue
		}

		// --- Parsing 6000 Series ---
		if matches := re6000ONU.FindStringSubmatch(line); matches != nil {
			slot, _ := strconv.Atoi(matches[1])
			port, _ := strconv.Atoi(matches[2])
			onuID, _ := strconv.Atoi(matches[3])
			model := matches[4]
			sn := matches[5]

			key := makeKey(slot, port, onuID)
			onuMap[key] = &models.ONU{
				SlotID:       slot,
				PortID:       port,
				OnuID:        onuID,
				SerialNumber: sn,
				Model:        model,
			}
			continue
		}

		// Mantém o estado da interface atual para a série 6000
		if matches := re6000Interface.FindStringSubmatch(line); matches != nil {
			currentSlot, _ = strconv.Atoi(matches[1])
			currentPort, _ = strconv.Atoi(matches[2])
			continue
		}

		if currentSlot > 0 && currentPort > 0 {
			if matches := re6000WanCfg.FindStringSubmatch(line); matches != nil {
				onuID, _ := strconv.Atoi(matches[1])
				user := matches[2]
				pass := matches[3]

				key := makeKey(currentSlot, currentPort, onuID)
				if onu, exists := onuMap[key]; exists {
					onu.PPPoEUser = user
					onu.PPPoEPass = pass
				}
				continue
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Converter o map para slice
	var onus []models.ONU
	for _, onu := range onuMap {
		onus = append(onus, *onu)
	}

	return onus, nil
}
