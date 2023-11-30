package esios

type EsioPVPC struct {
	Day  string `json:"Dia"`
	Hour string `json:"Hora"`
	PCB  string `json:"PCB"`
	GEN  string `json:"GEN"`
}

type EsiosResponse struct {
	PVPC    []EsioPVPC `json:"PVPC"`
	Message string     `json:"message"`
}
