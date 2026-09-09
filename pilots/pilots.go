package pilots

import (
	"os"

	"github.com/bomkz/cloudsync/pilots/vtscfg"
)

func ReadPilotSaveFile() PilotsFile {
	pf, err := os.ReadFile(`c:\Users\bomkz\AppData\Roaming\Boundless Dynamics, LLC\VTOLVR\SaveData\pilots.cfg`)
	if err != nil {
		panic(err)
	}

	var pilotFile PilotsFile
	err = vtscfg.Unmarshal(pf, &pilotFile)
	if err != nil {
		panic(err)
	}
	return pilotFile
}

func ReadGameSettingsFile() GameSettings {
	gf, err := os.ReadFile(`c:\Users\bomkz\AppData\Roaming\Boundless Dynamics, LLC\VTOLVR\SaveData\gameSettings.cfg`)
	if err != nil {
		panic(err)
	}

	var settingsFile GameSettings
	err = vtscfg.Unmarshal(gf, &settingsFile)
	if err != nil {
		panic(err)
	}

	return settingsFile
}
