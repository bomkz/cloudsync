package pilots

import (
	"os"

	"github.com/bomkz/cloudsync/global"
	"github.com/bomkz/cloudsync/pilots/vtscfg"
)

func ReadPilotSaveFile() global.PilotsFile {
	pf, err := os.ReadFile(`c:\Users\bomkz\AppData\Roaming\Boundless Dynamics, LLC\VTOLVR\SaveData\pilots.cfg`)
	if err != nil {
		panic(err)
	}

	var pilotFile global.PilotsFile
	err = vtscfg.Unmarshal(pf, &pilotFile)
	if err != nil {
		panic(err)
	}
	return pilotFile
}

func ReadGameSettingsFile() global.GameSettings {
	gf, err := os.ReadFile(`c:\Users\bomkz\AppData\Roaming\Boundless Dynamics, LLC\VTOLVR\SaveData\gameSettings.cfg`)
	if err != nil {
		panic(err)
	}

	var settingsFile global.GameSettings
	err = vtscfg.Unmarshal(gf, &settingsFile)
	if err != nil {
		panic(err)
	}

	return settingsFile
}
