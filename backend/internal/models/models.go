package models

type OLT struct {
	Vendor string
	Model  string
	Boards []Board
}

type Board struct {
	SlotID   int
	Type     string
	PONPorts []PONPort
}

type PONPort struct {
	PortID int
	ONUs   []ONU
}

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

type Service struct {
	ServiceID int
	Type      string
	CVLAN     int
	SVLAN     int
}

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
