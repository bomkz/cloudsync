package cloudsrv

import "github.com/bomkz/cloudsync/pilots"

type UpdatePilotReqStruct struct {
	Req   string            `json:"req"`
	Pilot pilots.PilotsFile `json:"body"`
	Auth  string            `json:"auth"`
}

type UpdatePilotPilotStruct struct {
	Name  string            `json:"name"`
	Pilot pilots.PilotsFile `json:"pilotFile"`
}

type GetPilotReqStruct struct {
	Req  string `json:"req"`
	Name string `json:"body"`
	Auth string `json:"auth"`
}

type GetPilotRespStruct struct {
	Req   string            `json:"req"`
	Pilot pilots.PilotsFile `json:"body"`
}

type OkRespStruct struct {
	Req string `json:"req"`
}

type ErrRespStruct struct {
	Req   string `json:"req"`
	Error string `json:"error"`
}
