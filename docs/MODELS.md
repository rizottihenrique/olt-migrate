# Modelos de Dados (Golang)

Todos os modelos de transição de estado na aplicação estão centralizados em `internal/models/models.go`.

### Estrutura ONU
Objeto central da aplicação, usado como "lingua franca" entre o `Parser` e o `Generator`.

```go
type ONU struct {
	SlotID       int
	PortID       int
	OnuID        int
	SerialNumber string
	Model        string
	PPPoEUser    string
	PPPoEPass    string
	Services     []Service
}
```

### Estrutura Service
Utilizada para agrupar as VLANs (Internet, VoIP, IPTV) de um determinado cliente na ONU.

```go
type Service struct {
	ServiceID int
	Type      string // ex: "INTERNET", "IPTV"
	CVLAN     int
	SVLAN     int
}
```

### Estruturas de Requisição (API)
Utilizadas pela camada Controller (`main.go`) e no pacote de `mapper`.

```go
type MigrationMapping struct {
	SourceSlot int
	SourcePON  int
	DestSlot   int
	DestPON    int
}

type MigrationRequest struct {
	ConfigFileContent []byte
	Mappings          []MigrationMapping
}
```
