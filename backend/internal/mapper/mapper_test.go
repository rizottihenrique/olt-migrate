package mapper

import (
	"testing"
	"olt-migrate-backend/internal/models"
)

func TestMapONUs(t *testing.T) {
	// 3 ONUs na OLT legada:
	// - ONU 1: Slot 1, PON 1
	// - ONU 2: Slot 1, PON 1
	// - ONU 3: Slot 2, PON 1
	original := []models.ONU{
		{SlotID: 1, PortID: 1, OnuID: 1, SerialNumber: "MAC1"},
		{SlotID: 1, PortID: 1, OnuID: 2, SerialNumber: "MAC2"},
		{SlotID: 2, PortID: 1, OnuID: 1, SerialNumber: "MAC3"},
	}

	// Regras:
	// Mover tudo do Slot 1/PON 1 -> Slot 3/PON 1 (Nokia)
	// Mover tudo do Slot 2/PON 1 -> Slot 3/PON 1 (Nokia)
	// Isso força um aglutinamento de duas PONs velhas em uma PON nova.
	rules := []models.MigrationMapping{
		{SourceSlot: 1, SourcePON: 1, DestSlot: 3, DestPON: 1},
		{SourceSlot: 2, SourcePON: 1, DestSlot: 3, DestPON: 1},
	}

	mapped := MapONUs(original, rules)

	if len(mapped) != 3 {
		t.Fatalf("Esperava 3 ONUs mapeadas, obteve %d", len(mapped))
	}

	// Todas devem estar no Slot 3, PON 1
	for _, onu := range mapped {
		if onu.SlotID != 3 || onu.PortID != 1 {
			t.Errorf("ONU não migrada corretamente: %+v", onu)
		}
	}

	// IDs não podem colidir
	ids := make(map[int]bool)
	for _, onu := range mapped {
		if ids[onu.OnuID] {
			t.Errorf("Colisão de OnuID detectada: ID %d repetido", onu.OnuID)
		}
		ids[onu.OnuID] = true
	}
}
