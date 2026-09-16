package main

import "github.com/bomkz/cloudsync/pilots"

func main() {
	pilotsfile := pilots.ReadPilotSaveFile()

	updatePilotReq := cloudstore.updatePilotRequestStruct{}
}
