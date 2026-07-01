package mapper

import (
	"fmt"
	"olt-migrate-backend/internal/models"
)

// MapONUs recebe a lista original de ONUs e as regras de mapeamento,
// e retorna uma nova lista de ONUs transpostas para os novos slots e portas.
// Também reatribui os IDs das ONUs para evitar conflitos se múltiplas portas PON
// de origem forem aglutinadas em uma mesma porta PON de destino.
func MapONUs(originalONUs []models.ONU, mappings []models.MigrationMapping) []models.ONU {
	var mappedONUs []models.ONU

	// Mapeamento rápido de regras: chave "SourceSlot-SourcePON" -> mapping
	ruleMap := make(map[string]models.MigrationMapping)
	for _, m := range mappings {
		key := fmt.Sprintf("%d-%d", m.SourceSlot, m.SourcePON)
		ruleMap[key] = m
	}

	// Contador de OnuIDs livres por porta PON de destino (chave: "DestSlot-DestPON")
	destOnuCounters := make(map[string]int)

	for _, onu := range originalONUs {
		sourceKey := fmt.Sprintf("%d-%d", onu.SlotID, onu.PortID)

		// Se existir uma regra de mapeamento para esta porta
		if rule, exists := ruleMap[sourceKey]; exists {
			// Clona a ONU
			mappedONU := onu

			// Atualiza Slot e PON
			mappedONU.SlotID = rule.DestSlot
			mappedONU.PortID = rule.DestPON

			// Resolve OnuID para evitar conflito
			destKey := fmt.Sprintf("%d-%d", rule.DestSlot, rule.DestPON)
			
			// Incrementa o contador para achar o próximo ID livre (inicia em 1)
			destOnuCounters[destKey]++
			mappedONU.OnuID = destOnuCounters[destKey]

			mappedONUs = append(mappedONUs, mappedONU)
		} else {
			// Se não houver regra, podemos optar por não migrar ou manter a original.
			// Na nossa regra de negócio, só migramos o que está na tabela de de/para.
			// Portanto, ignoramos as ONUs que não têm mapeamento de porta.
		}
	}

	return mappedONUs
}
