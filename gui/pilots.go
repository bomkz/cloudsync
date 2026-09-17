package gui

import (
	"fmt"
	"image/color"
	"log"
	"sort"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/bomkz/cloudsync/cloudsrv"
	"github.com/bomkz/cloudsync/global"
	"github.com/bomkz/cloudsync/pilots"
	"github.com/bomkz/cloudsync/pilots/vtscfg"
	"golang.org/x/image/colornames"
)

func buildPilotsTabGui() *container.TabItem {
	pilotsS, err := cloudsrv.GetPilots()
	if err != nil {
		log.Fatal(err)
	}

	pilotSaveList := widget.NewSelect(pilotsS, nil)
	var versionsInt int
	if pilotSaveList.Selected != "" {
		versionsInt, err = cloudsrv.GetPilotInfo(pilotSaveList.Selected)
		if err != nil {
			log.Fatal(err)
		}

	}
	var versions []string

	for x := range versionsInt - 1 {
		versions = append(versions, fmt.Sprint(x+1))
	}

	pilotSaveVariantList := widget.NewSelect(versions, nil)

	saveButton := widget.NewButton("Save", nil)
	delButton := widget.NewButton("Delete", nil)
	loadButton := widget.NewButton("Load", nil)
	newButton := widget.NewButton("New", nil)

	newButtonFunc := func() {
		nameEntry := widget.NewEntry()
		nameInput := widget.NewFormItem("Name", nameEntry)

		var newSave *dialog.FormDialog
		newSave = dialog.NewForm("Create New Pilot Save", "Save", "Cancel", []*widget.FormItem{nameInput}, func(b bool) {
			if !b {
				return
			}

			name := strings.TrimSpace(nameEntry.Text)
			if name == "" {
				newSave.Hide()
				dialog.NewError(fmt.Errorf("Pilot name cannot be empty!"), w).Show()
				return
			}
			if !isAlphanumericalName(name) {
				newSave.Hide()
				dialog.NewError(fmt.Errorf("Pilot name must contain alphanumerical characters only!"), w).Show()
				return
			}

			for _, existingName := range pilotsS {
				if strings.EqualFold(existingName, name) {
					newSave.Hide()
					dialog.NewError(fmt.Errorf("A pilot named %q already exists!", existingName), w).Show()
					return
				}
			}

			currentPilot := pilots.ReadPilotSaveFile()
			err = cloudsrv.UpdatePilot(name, currentPilot)
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}

			pilotsS, err = cloudsrv.GetPilots()
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			pilotSaveList.Options = pilotsS
			pilotSaveList.Refresh()
			pilotSaveList.SetSelected(name)

		}, w)
		newSave.Show()
	}
	newButton.OnTapped = newButtonFunc

	header := container.NewHBox(newButton, saveButton, delButton, loadButton, pilotSaveList, pilotSaveVariantList)
	currentData := buildPilotsTabCurrentDataGui()

	if pilotSaveVariantList.Selected == "" {
		// The saved-data container is populated below after it is mounted.
	} else {
		// The selected version is loaded by the change handler below.
	}

	if pilotSaveList.Selected == "" {
		saveButton.Disable()
		delButton.Disable()
		loadButton.Disable()
	}

	savedPilotContainer := container.NewStack()
	content := container.NewHSplit(currentData, savedPilotContainer)
	content.SetOffset(0.5)
	pilotContainer := container.NewTabItem("Pilot Sync", container.NewBorder(header, nil, nil, nil, content))

	setSavedPilotData := func(pilotData global.PilotsFile) {
		savedPilotContainer.RemoveAll()
		savedPilotContainer.Add(buildPilotsTabSavedDataGui(pilotData))
		savedPilotContainer.Refresh()
	}
	setSavedPilotData(global.PilotsFile{})

	onChangedPilotSave := func(s string) {
		if s == "" {
			saveButton.Disable()
			delButton.Disable()
			loadButton.Disable()
			setSavedPilotData(global.PilotsFile{})
			pilotSaveVariantList.Options = nil
			pilotSaveVariantList.ClearSelected()
			pilotSaveVariantList.Refresh()
			return
		}
		saveButton.Enable()
		delButton.Enable()
		loadButton.Enable()

		versionsInt, err = cloudsrv.GetPilotInfo(pilotSaveList.Selected)
		if err != nil {
			log.Fatal(err)
		}
		versions = versions[:0]
		for x := range versionsInt - 1 {
			versions = append(versions, fmt.Sprint(x+1))
		}
		pilotSaveVariantList.Options = versions
		pilotSaveVariantList.Refresh()
	}
	onChangedPilotVariantSave := func(s string) {
		if s == "" {
			setSavedPilotData(global.PilotsFile{})
			return
		}

		pilotData, err := cloudsrv.GetPilot(pilotSaveList.Selected, s)
		if err != nil {
			dialog.NewError(err, w).Show()
			return
		}
		setSavedPilotData(pilotData)
	}

	pilotSaveVariantList.OnChanged = onChangedPilotVariantSave

	pilotSaveList.OnChanged = onChangedPilotSave

	return pilotContainer
}

func isAlphanumericalName(name string) bool {
	for _, character := range name {
		if !unicode.IsLetter(character) && !unicode.IsNumber(character) {
			return false
		}
	}
	return name != ""
}

func buildPilotsTabSavedDataGui(pilotData global.PilotsFile) fyne.CanvasObject {

	savedPilotData := container.NewVBox()

	savedLabel := widget.NewLabel("Saved Pilot Data")
	savedPilotData.Add(savedLabel)

	if len(pilotData.Pilots.PilotSaves) == 0 {
		savedPilotData.Add(widget.NewLabel("No pilot data found."))
	}

	var pilotNames []string
	for _, pilot := range pilotData.Pilots.PilotSaves {
		pilotNames = append(pilotNames, pilot.PilotName)
	}

	pilotList := widget.NewSelect(pilotNames, nil)

	pilotName := widget.NewLabel(fmt.Sprintf("Pilot: %s", "N/A"))
	pilotVehicle := widget.NewLabel(fmt.Sprintf("Last vehicle: %s", "N/A"))
	pilotFlightTime := widget.NewLabel("Total flight time: N/A")

	gsuitColor, gsuitSwatch, gsuitColorWidget := newColorRow("G Suit Color")
	skinColor, skinSwatch, skinColorWidget := newColorRow("Skin Color")
	strapColor, strapSwatch, strapColorWidget := newColorRow("Strap Color")
	suitColor, suitSwatch, suitColorWidget := newColorRow("Suit Color")
	vestColor, vestSwatch, vestColorWidget := newColorRow("Vest Color")

	pilotColor := container.NewVBox(gsuitColorWidget, skinColorWidget, strapColorWidget, suitColorWidget, vestColorWidget)
	vehicleList := widget.NewSelect(nil, nil)
	currentVehicleData := container.NewVBox(
		widget.NewLabel("Vehicle"),
		vehicleList,
	)
	currentVehicleDetails := container.NewVBox()
	currentCampaignDetails := container.NewVBox()
	currentVehicleData.Add(currentVehicleDetails)
	campaignList := widget.NewSelect(nil, nil)
	currentVehicleData.Add(widget.NewLabel("Campaign vehicle configuration"))
	currentVehicleData.Add(campaignList)
	currentVehicleData.Add(currentCampaignDetails)

	updateCampaign := func(selected string, vehicle global.Vehicle) {
		currentCampaignDetails.RemoveAll()
		for _, campaign := range vehicle.Campaigns {
			campaignKey := campaign.CampaignName
			if campaignKey == "" {
				campaignKey = campaign.CampaignID
			}
			if campaignKey == selected {
				addCampaignDetails(currentCampaignDetails, campaign)
				currentCampaignDetails.Refresh()
				return
			}
		}
		currentCampaignDetails.Add(widget.NewLabel("No campaign selected."))
		currentCampaignDetails.Refresh()
	}

	updateVehicle := func(selected string, pilot global.PilotSave) {
		currentVehicleDetails.RemoveAll()
		currentCampaignDetails.RemoveAll()
		for _, vehicle := range pilot.Vehicles {
			if vehicle.VehicleName == selected {
				addVehicleDetails(currentVehicleDetails, vehicle)

				campaigns := make([]string, 0, len(vehicle.Campaigns))
				for _, campaign := range vehicle.Campaigns {
					campaignKey := campaign.CampaignName
					if campaignKey == "" {
						campaignKey = campaign.CampaignID
					}
					campaigns = append(campaigns, campaignKey)
				}
				campaignList.Options = campaigns
				campaignList.Refresh()
				if len(campaigns) > 0 {
					campaignList.SetSelected(campaigns[0])
				} else {
					currentCampaignDetails.Add(widget.NewLabel("No campaign data found."))
				}
				currentVehicleDetails.Refresh()
				return
			}
		}
		campaignList.Options = nil
		campaignList.ClearSelected()
		campaignList.Refresh()
		currentVehicleDetails.Add(widget.NewLabel("No vehicle selected."))
		currentVehicleDetails.Refresh()
	}

	updatePilot := func(selected string) {
		for _, pilot := range pilotData.Pilots.PilotSaves {
			if pilot.PilotName != selected {
				continue
			}

			pilotName.SetText(fmt.Sprintf("Pilot: %s", pilot.PilotName))
			pilotVehicle.SetText(fmt.Sprintf("Last vehicle: %s", pilot.LastVehicleUsed))
			pilotFlightTime.SetText(fmt.Sprintf("Total flight time: %.2f hours", (pilot.TotalFlightTime/60)/60))
			setColorRow(gsuitColor, gsuitSwatch, "G Suit Color", pilot.GSuitColor)
			setColorRow(skinColor, skinSwatch, "Skin Color", pilot.SkinColor)
			setColorRow(strapColor, strapSwatch, "Strap Color", pilot.StrapsColor)
			setColorRow(suitColor, suitSwatch, "Suit Color", pilot.SuitColor)
			setColorRow(vestColor, vestSwatch, "Vest Color", pilot.VestColor)

			vehicles := make([]string, 0, len(pilot.Vehicles))
			for _, vehicle := range pilot.Vehicles {
				vehicles = append(vehicles, vehicle.VehicleName)
			}
			vehicleList.Options = vehicles
			vehicleList.Refresh()
			if len(vehicles) > 0 {
				vehicleList.SetSelected(vehicles[0])
			} else {
				vehicleList.ClearSelected()
				campaignList.Options = nil
				campaignList.ClearSelected()
				campaignList.Refresh()
				currentVehicleDetails.RemoveAll()
				currentVehicleDetails.Add(widget.NewLabel("No vehicle data found."))
				currentCampaignDetails.RemoveAll()
				currentCampaignDetails.Add(widget.NewLabel("No campaign data found."))
				currentVehicleDetails.Refresh()
				currentCampaignDetails.Refresh()
			}

			return
		}
	}

	pilotList.OnChanged = updatePilot

	savedPilotData.Add(pilotList)
	savedPilotData.Add(pilotName)
	savedPilotData.Add(pilotVehicle)
	savedPilotData.Add(pilotFlightTime)
	savedPilotData.Add(pilotColor)

	vehicleList.OnChanged = func(selected string) {
		for _, pilot := range pilotData.Pilots.PilotSaves {
			if pilot.PilotName != pilotList.Selected {
				continue
			}
			updateVehicle(selected, pilot)
			return
		}
	}
	campaignList.OnChanged = func(selected string) {
		for _, pilot := range pilotData.Pilots.PilotSaves {
			if pilot.PilotName != pilotList.Selected {
				continue
			}
			for _, vehicle := range pilot.Vehicles {
				if vehicle.VehicleName == vehicleList.Selected {
					updateCampaign(selected, vehicle)
					return
				}
			}
		}
	}

	currentData := container.NewVBox(savedPilotData, currentVehicleData)
	if len(pilotNames) > 0 {
		pilotList.SetSelected(pilotNames[0])
	} else {
		currentVehicleDetails.Add(widget.NewLabel("No vehicle data found."))
	}

	return container.NewVScroll(currentData)

}

func buildPilotsTabCurrentDataGui() fyne.CanvasObject {
	currentPilotData := container.NewVBox()

	currentLabel := widget.NewLabel("Current Pilot Data")
	currentPilotData.Add(currentLabel)
	pilotData := pilots.ReadPilotSaveFile()
	if len(pilotData.Pilots.PilotSaves) == 0 {
		currentPilotData.Add(widget.NewLabel("No pilot data found."))
	}

	var pilotNames []string
	for _, pilot := range pilotData.Pilots.PilotSaves {
		pilotNames = append(pilotNames, pilot.PilotName)
	}

	pilotList := widget.NewSelect(pilotNames, nil)

	pilotName := widget.NewLabel(fmt.Sprintf("Pilot: %s", "N/A"))
	pilotVehicle := widget.NewLabel(fmt.Sprintf("Last vehicle: %s", "N/A"))
	pilotFlightTime := widget.NewLabel("Total flight time: N/A")

	gsuitColor, gsuitSwatch, gsuitColorWidget := newColorRow("G Suit Color")
	skinColor, skinSwatch, skinColorWidget := newColorRow("Skin Color")
	strapColor, strapSwatch, strapColorWidget := newColorRow("Strap Color")
	suitColor, suitSwatch, suitColorWidget := newColorRow("Suit Color")
	vestColor, vestSwatch, vestColorWidget := newColorRow("Vest Color")

	pilotColor := container.NewVBox(gsuitColorWidget, skinColorWidget, strapColorWidget, suitColorWidget, vestColorWidget)
	vehicleList := widget.NewSelect(nil, nil)
	currentVehicleData := container.NewVBox(
		widget.NewLabel("Vehicle"),
		vehicleList,
	)
	currentVehicleDetails := container.NewVBox()
	currentCampaignDetails := container.NewVBox()
	currentVehicleData.Add(currentVehicleDetails)
	campaignList := widget.NewSelect(nil, nil)
	currentVehicleData.Add(widget.NewLabel("Campaign vehicle configuration"))
	currentVehicleData.Add(campaignList)
	currentVehicleData.Add(currentCampaignDetails)

	updateCampaign := func(selected string, vehicle global.Vehicle) {
		currentCampaignDetails.RemoveAll()
		for _, campaign := range vehicle.Campaigns {
			campaignKey := campaign.CampaignName
			if campaignKey == "" {
				campaignKey = campaign.CampaignID
			}
			if campaignKey == selected {
				addCampaignDetails(currentCampaignDetails, campaign)
				currentCampaignDetails.Refresh()
				return
			}
		}
		currentCampaignDetails.Add(widget.NewLabel("No campaign selected."))
		currentCampaignDetails.Refresh()
	}

	updateVehicle := func(selected string, pilot global.PilotSave) {
		currentVehicleDetails.RemoveAll()
		currentCampaignDetails.RemoveAll()
		for _, vehicle := range pilot.Vehicles {
			if vehicle.VehicleName == selected {
				addVehicleDetails(currentVehicleDetails, vehicle)

				campaigns := make([]string, 0, len(vehicle.Campaigns))
				for _, campaign := range vehicle.Campaigns {
					campaignKey := campaign.CampaignName
					if campaignKey == "" {
						campaignKey = campaign.CampaignID
					}
					campaigns = append(campaigns, campaignKey)
				}
				campaignList.Options = campaigns
				campaignList.Refresh()
				if len(campaigns) > 0 {
					campaignList.SetSelected(campaigns[0])
				} else {
					currentCampaignDetails.Add(widget.NewLabel("No campaign data found."))
				}
				currentVehicleDetails.Refresh()
				return
			}
		}
		campaignList.Options = nil
		campaignList.ClearSelected()
		campaignList.Refresh()
		currentVehicleDetails.Add(widget.NewLabel("No vehicle selected."))
		currentVehicleDetails.Refresh()
	}

	updatePilot := func(selected string) {
		for _, pilot := range pilotData.Pilots.PilotSaves {
			if pilot.PilotName != selected {
				continue
			}

			pilotName.SetText(fmt.Sprintf("Pilot: %s", pilot.PilotName))
			pilotVehicle.SetText(fmt.Sprintf("Last vehicle: %s", pilot.LastVehicleUsed))
			pilotFlightTime.SetText(fmt.Sprintf("Total flight time: %.2f hours", (pilot.TotalFlightTime/60)/60))
			setColorRow(gsuitColor, gsuitSwatch, "G Suit Color", pilot.GSuitColor)
			setColorRow(skinColor, skinSwatch, "Skin Color", pilot.SkinColor)
			setColorRow(strapColor, strapSwatch, "Strap Color", pilot.StrapsColor)
			setColorRow(suitColor, suitSwatch, "Suit Color", pilot.SuitColor)
			setColorRow(vestColor, vestSwatch, "Vest Color", pilot.VestColor)

			vehicles := make([]string, 0, len(pilot.Vehicles))
			for _, vehicle := range pilot.Vehicles {
				vehicles = append(vehicles, vehicle.VehicleName)
			}
			vehicleList.Options = vehicles
			vehicleList.Refresh()
			if len(vehicles) > 0 {
				vehicleList.SetSelected(vehicles[0])
			} else {
				vehicleList.ClearSelected()
				campaignList.Options = nil
				campaignList.ClearSelected()
				campaignList.Refresh()
				currentVehicleDetails.RemoveAll()
				currentVehicleDetails.Add(widget.NewLabel("No vehicle data found."))
				currentCampaignDetails.RemoveAll()
				currentCampaignDetails.Add(widget.NewLabel("No campaign data found."))
				currentVehicleDetails.Refresh()
				currentCampaignDetails.Refresh()
			}

			return
		}
	}

	pilotList.OnChanged = updatePilot

	currentPilotData.Add(pilotList)
	currentPilotData.Add(pilotName)
	currentPilotData.Add(pilotVehicle)
	currentPilotData.Add(pilotFlightTime)
	currentPilotData.Add(pilotColor)

	vehicleList.OnChanged = func(selected string) {
		for _, pilot := range pilotData.Pilots.PilotSaves {
			if pilot.PilotName != pilotList.Selected {
				continue
			}
			updateVehicle(selected, pilot)
			return
		}
	}
	campaignList.OnChanged = func(selected string) {
		for _, pilot := range pilotData.Pilots.PilotSaves {
			if pilot.PilotName != pilotList.Selected {
				continue
			}
			for _, vehicle := range pilot.Vehicles {
				if vehicle.VehicleName == vehicleList.Selected {
					updateCampaign(selected, vehicle)
					return
				}
			}
		}
	}

	currentData := container.NewVBox(currentPilotData, currentVehicleData)
	if len(pilotNames) > 0 {
		pilotList.SetSelected(pilotNames[0])
	} else {
		currentVehicleDetails.Add(widget.NewLabel("No vehicle data found."))
	}

	return container.NewVScroll(currentData)
}

func addVehicleDetails(details *fyne.Container, vehicle global.Vehicle) {
	details.Add(widget.NewLabel(fmt.Sprintf("Vehicle: %s", vehicle.VehicleName)))
	details.Add(widget.NewLabel(fmt.Sprintf("Seat height: %.3f", vehicle.SeatHeight)))
	details.Add(widget.NewLabel(fmt.Sprintf("Joystick position: (%.3f, %.3f, %.3f)", vehicle.JoystickPosition[0], vehicle.JoystickPosition[1], vehicle.JoystickPosition[2])))
	details.Add(widget.NewLabel(fmt.Sprintf("Throttle position: (%.3f, %.3f, %.3f)", vehicle.ThrottlePosition[0], vehicle.ThrottlePosition[1], vehicle.ThrottlePosition[2])))
	details.Add(widget.NewLabel(fmt.Sprintf("Altitude mode: %s", vehicle.AltitudeMode)))
	details.Add(widget.NewLabel(fmt.Sprintf("Distance mode: %s", vehicle.DistanceMode)))
	details.Add(widget.NewLabel(fmt.Sprintf("Airspeed mode: %s", vehicle.AirspeedMode)))
	details.Add(widget.NewLabel(fmt.Sprintf("Last livery ID: %s", vehicle.LastLiveryID)))

	details.Add(widget.NewLabel("VDATA"))
	addNodeDetails(details, vehicle.VData, 1)

	details.Add(widget.NewLabel("Saved loadouts"))
	for _, loadout := range vehicle.SavedLoadouts.Loadouts {
		details.Add(widget.NewLabel(fmt.Sprintf("Loadout: %s | Fuel: %.3f", loadout.Name, loadout.NormFuel)))
		for index, weapon := range []string{loadout.Eq0, loadout.Eq1, loadout.Eq2, loadout.Eq3, loadout.Eq4, loadout.Eq5, loadout.Eq6, loadout.Eq7, loadout.Eq8, loadout.Eq9, loadout.Eq10, loadout.Eq11, loadout.Eq12, loadout.Eq13, loadout.Eq14, loadout.Eq15} {
			if weapon != "" {
				details.Add(widget.NewLabel(fmt.Sprintf("  Equipment %d: %s", index, weapon)))
			}
		}
	}
	details.Refresh()
}

func addCampaignDetails(details *fyne.Container, campaign global.Campaign) {
	name := campaign.CampaignName
	if name == "" {
		name = "Unnamed campaign"
	}
	details.Add(widget.NewLabel(fmt.Sprintf("Campaign: %s", name)))
	details.Add(widget.NewLabel(fmt.Sprintf("Campaign ID: %s", campaign.CampaignID)))
	details.Add(widget.NewLabel(fmt.Sprintf("Vehicle: %s", campaign.VehicleName)))
	details.Add(widget.NewLabel(fmt.Sprintf("Current fuel: %.3f", campaign.CurrentFuel)))
	details.Add(widget.NewLabel(fmt.Sprintf("Last scenario index: %d", campaign.LastScenarioIndex)))
	details.Add(widget.NewLabel(fmt.Sprintf("Last scenario was training: %t", campaign.LastScenarioTraining)))
	details.Add(widget.NewLabel("Current weapons"))
	for _, weapon := range campaign.CurrentWeapons.Slots {
		details.Add(widget.NewLabel(fmt.Sprintf("Weapon %d: %s", weapon.Index, weapon.Weapon)))
	}
	details.Refresh()
}

func addNodeDetails(details *fyne.Container, node vtscfg.Node, depth int) {
	for _, field := range node.Fields {
		indent := ""
		for index := 0; index < depth; index++ {
			indent += "  "
		}
		if field.IsScalar() {
			details.Add(widget.NewLabel(fmt.Sprintf("%s%s: %s", indent, field.Key, field.RawValue())))
			continue
		}
		if field.IsNode() {
			details.Add(widget.NewLabel(fmt.Sprintf("%s%s", indent, field.Key)))
			addNodeDetails(details, *field.Child, depth+1)
		}
	}
}

func newColorRow(name string) (*widget.Label, *canvas.Rectangle, *fyne.Container) {
	label := widget.NewLabel(name)
	swatch := canvas.NewRectangle(color.NRGBA{A: 255})
	swatch.SetMinSize(fyne.NewSize(24, 18))
	return label, swatch, container.NewHBox(swatch, label)
}

func setColorRow(label *widget.Label, swatch *canvas.Rectangle, name string, rgb [3]float64) {
	swatchColor := colorFromRGB(rgb)
	label.SetText(fmt.Sprintf("%s: #%02X%02X%02X (%s)", name, swatchColor.R, swatchColor.G, swatchColor.B, nearestColorName(swatchColor)))
	swatch.FillColor = swatchColor
	swatch.Refresh()
}

func colorFromRGB(rgb [3]float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(clampColor(rgb[0]) * 255),
		G: uint8(clampColor(rgb[1]) * 255),
		B: uint8(clampColor(rgb[2]) * 255),
		A: 255,
	}
}

func nearestColorName(value color.NRGBA) string {
	bestName := ""
	bestDistance := -1
	for _, name := range cssColorNames {
		candidate := colornames.Map[name]
		redDistance := int(value.R) - int(candidate.R)
		greenDistance := int(value.G) - int(candidate.G)
		blueDistance := int(value.B) - int(candidate.B)
		distance := redDistance*redDistance + greenDistance*greenDistance + blueDistance*blueDistance
		if bestDistance == -1 || distance < bestDistance {
			bestName = name
			bestDistance = distance
		}
	}
	return bestName
}

var cssColorNames = func() []string {
	names := make([]string, 0, len(colornames.Map))
	for name := range colornames.Map {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}()

func clampColor(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

// DISCLAIMER: This file has been partially LLM Generated.
