package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"olt-migrate-backend/internal/core"
	"olt-migrate-backend/internal/mapper"
	"olt-migrate-backend/internal/models"
)

func migrateHandler(w http.ResponseWriter, r *http.Request) {
	// CORS handling para chamadas do Frontend Next.js (dev mode)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Faz o parse do multipart form (upload do arquivo)
	err := r.ParseMultipartForm(50 << 20) // 50 MB
	if err != nil {
		http.Error(w, "Falha ao ler formulário multipart", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("configFile")
	if err != nil {
		http.Error(w, "Arquivo de configuração é obrigatório", http.StatusBadRequest)
		return
	}
	defer file.Close()

	sourceVendor := r.FormValue("sourceVendor")
	destVendor := r.FormValue("destVendor")

	// Lê as regras de mapeamento do JSON enviado no formulário
	mappingsJSON := r.FormValue("mappings")
	var mappings []models.MigrationMapping
	if mappingsJSON != "" {
		if err := json.Unmarshal([]byte(mappingsJSON), &mappings); err != nil {
			http.Error(w, "Erro ao decodificar as regras de mapeamento", http.StatusBadRequest)
			return
		}
	}

	// 1. Extração: Obtém o Parser dinamicamente
	parser := core.GetParser(sourceVendor)
	onus, err := parser.Parse(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao processar arquivo: %v", err), http.StatusInternalServerError)
		return
	}

	// 2. Mapeamento
	if len(mappings) > 0 {
		onus = mapper.MapONUs(onus, mappings)
	}

	// 3. Geração: Obtém o Generator dinamicamente
	generator := core.GetGenerator(destVendor)
	commands := generator.Generate(onus)

	// Agrupa ONUs por PON para gerar comandos modulares por porta (cópia rápida na UI)
	ponCommands := make(map[string]string)
	ponMap := make(map[string][]models.ONU)
	onuCommands := make(map[string]string)

	for _, onu := range onus {
		ponKey := fmt.Sprintf("PON %d/%d", onu.SlotID, onu.PortID)
		ponMap[ponKey] = append(ponMap[ponKey], onu)

		clientName := onu.PPPoEUser
		if clientName == "" {
			clientName = fmt.Sprintf("Bridge (SN: %s)", onu.SerialNumber)
		}
		onuKey := fmt.Sprintf("%d/%d/%d - %s", onu.SlotID, onu.PortID, onu.OnuID, clientName)
		onuCommands[onuKey] = generator.Generate([]models.ONU{onu})
	}
	for key, ponONUs := range ponMap {
		ponCommands[key] = generator.Generate(ponONUs)
	}

	// Retorna resultado JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"totalONUsFound": len(onus),
		"commands":       commands,
		"ponCommands":    ponCommands,
		"onuCommands":    onuCommands,
		"onus":           onus,
	})
}

func main() {
	http.HandleFunc("/api/migrate", migrateHandler)
	
	port := "8080"
	fmt.Printf("OLT Migrate Backend rodando na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
