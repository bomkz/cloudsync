package gui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func InitGui() {
	a = app.New()
	w = a.NewWindow("CloudSync.")

	buildGui()

	w.ShowAndRun()
}

func buildGui() {
	tabs := container.NewAppTabs()
	pilotTab := buildPilotsTabGui()
	configTab := container.NewTabItem("Config Sync", widget.NewLabel("Config Sync"))

	tabs.Append(pilotTab)
	tabs.Append(configTab)

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
}
