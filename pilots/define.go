package pilots

import (
	"time"

	"github.com/bomkz/cloudsync/pilots/vtscfg"
)

type UpdatePilotRequestStruct struct {
	Request string                       `json:"request"`
	Body    UpdatePilotRequestBodyStruct `json:"body"`
}

type UpdatePilotRequestBodyStruct struct {
	Token     string     `json:"token"`
	Name      string     `json:"name"`
	PilotData PilotsFile `json:"pilotData"`
}

type PilotsFile struct {
	Pilots Pilots `vts:"PILOTS" json:"PILOTS"`
}

type Pilots struct {
	PilotSaves []PilotSave `vts:"PILOTSAVE" json:"PILOTSAVE"`
}

type PilotSave struct {
	PilotName       string     `vts:"pilotName" json:"pilotName"`
	LastVehicleUsed string     `vts:"lastVehicleUsed" json:"lastVehicleUsed"`
	TotalFlightTime float64    `vts:"totalFlightTime" json:"totalFlightTime"`
	SkinColor       [3]float64 `vts:"skinColor" json:"skinColor"`
	SuitColor       [3]float64 `vts:"suitColor" json:"suitColor"`
	VestColor       [3]float64 `vts:"vestColor" json:"vestColor"`
	GSuitColor      [3]float64 `vts:"gSuitColor" json:"gSuitColor"`
	StrapsColor     [3]float64 `vts:"strapsColor" json:"strapsColor"`
	Vehicles        []Vehicle  `vts:"VEHICLE" json:"VEHICLE"`
}

type Vehicle struct {
	VehicleName      string        `vts:"vehicleName" json:"vehicleName"`
	SeatHeight       float64       `vts:"seatHeight" json:"seatHeight"`
	JoystickPosition [3]float64    `vts:"joystickPosition" json:"joystickPosition"`
	ThrottlePosition [3]float64    `vts:"throttlePosition" json:"throttlePosition"`
	AltitudeMode     string        `vts:"altitudeMode" json:"altitudeMode"`
	DistanceMode     string        `vts:"distanceMode" json:"distanceMode"`
	AirspeedMode     string        `vts:"airspeedMode" json:"airspeedMode"`
	LastLiveryID     string        `vts:"lastLiveryId" json:"lastLiveryId"`
	VData            vtscfg.Node   `vts:"VDATA" json:"VDATA"`
	Campaigns        []Campaign    `vts:"CAMPAIGN" json:"CAMPAIGN"`
	SavedLoadouts    SavedLoadouts `vts:"SavedLoadouts" json:"SavedLoadouts"`
}

type Campaign struct {
	CampaignName         string         `vts:"campaignName" json:"campaignName"`
	CampaignID           string         `vts:"campaignID" json:"campaignID"`
	VehicleName          string         `vts:"vehicleName" json:"vehicleName"`
	AvailableWeapons     string         `vts:"availableWeapons" json:"availableWeapons"`
	CurrentFuel          float64        `vts:"currentFuel" json:"currentFuel"`
	AvailableScenarios   string         `vts:"availableScenarios" json:"availableScenarios"`
	LastScenarioIndex    int            `vts:"lastScenarioIdx" json:"lastScenarioIdx"`
	LastScenarioTraining bool           `vts:"lastScenarioWasTraining" json:"lastScenarioWasTraining"`
	CurrentWeapons       CurrentWeapons `vts:"currentWeapons" json:"currentWeapons"`
}

type CurrentWeapons struct {
	Slots []WeaponSlot `vts:"weapon" json:"weapon"`
}

type WeaponSlot struct {
	Index  int    `vts:"idx" json:"idx"`
	Weapon string `vts:"weapon" json:"weapon"`
}

type SavedLoadouts struct {
	Loadouts []SavedLoadout `vts:"SavedLoadout" json:"SavedLoadout"`
}

type SavedLoadout struct {
	Name     string  `vts:"name" json:"name"`
	NormFuel float64 `vts:"normFuel" json:"normFuel"`
	Eq0      string  `vts:"eq0" json:"eq0"`
	Eq1      string  `vts:"eq1" json:"eq1"`
	Eq2      string  `vts:"eq2" json:"eq2"`
	Eq3      string  `vts:"eq3" json:"eq3"`
	Eq4      string  `vts:"eq4" json:"eq4"`
	Eq5      string  `vts:"eq5" json:"eq5"`
	Eq6      string  `vts:"eq6" json:"eq6"`
	Eq7      string  `vts:"eq7" json:"eq7"`
	Eq8      string  `vts:"eq8" json:"eq8"`
	Eq9      string  `vts:"eq9" json:"eq9"`
	Eq10     string  `vts:"eq10" json:"eq10"`
	Eq11     string  `vts:"eq11" json:"eq11"`
	Eq12     string  `vts:"eq12" json:"eq12"`
	Eq13     string  `vts:"eq13" json:"eq13"`
	Eq14     string  `vts:"eq14" json:"eq14"`
	Eq15     string  `vts:"eq15" json:"eq15"`
}

type GameSettings struct {
	Settings GameSettingsWrapper `vts:"GAMESETTINGS" json:"GAMESETTINGS"`
}
type GameSettingsWrapper struct {
	RadioMusicPath        string  `vts:"RADIO_MUSIC_PATH" json:"RADIO_MUSIC_PATH"`
	WingmanVoices         string  `vts:"WINGMAN_VOICES" json:"WINGMAN_VOICES"`
	ToolTips              bool    `vts:"TOOLTIPS" json:"TOOLTIPS"`
	UnitIcons             bool    `vts:"UNIT_ICONS" json:"UNIT_ICONS"`
	BodyPhysics           bool    `vts:"BODY_PHYSICS" json:"BODY_PHYSICS"`
	TreeCollisions        bool    `vts:"TREE_COLLISIONS" json:"TREE_COLLISIONS"`
	HookPhysics           bool    `vts:"HOOK_PHYSICS" json:"HOOK_PHYSICS"`
	PersistentSCam        bool    `vts:"PERSISTENT_S_CAM" json:"PERSISTENT_S_CAM"`
	ShowBobbleHead        bool    `vts:"SHOW_BOBBLEHEAD" json:"SHOW_BOBBLEHEAD"`
	ThumbstickMode        bool    `vts:"THUMBSTICK_MODE" json:"THUMBSTICK_MODE"`
	ThumbstickDeadZone    int     `vts:"THUMBSTICK_DEADZONE" json:"THUMBSTICK_DEADZONE"`
	ThumbRudder           bool    `vts:"THUMB_RUDDER" json:"THUMB_RUDDER"`
	HardwareRudder        bool    `vts:"HARDWARE_RUDDER" json:"HARDWARE_RUDDER"`
	TapToggleGrip         bool    `vts:"TAP_TOGGLE_GRIP" json:"TAP_TOGGLE_GRIP"`
	ControlHaptics        int     `vts:"CONTROL_HAPTICS" json:"CONTROL_HAPTICS"`
	OverallHaptics        int     `vts:"OVERALL_HAPTICS" json:"OVERALL_HAPTICS"`
	Msaa                  int     `vts:"MSAA" json:"MSAA"`
	HideHelmet            bool    `vts:"HIDE_HELMET" json:"HIDE_HELMET"`
	FullScreenNVG         bool    `vts:"FULLSCREEN_NVG" json:"FULLSCREEN_NVG"`
	NvgPhosphor           bool    `vts:"NVG_PHOSPHOR" json:"NVG_PHOSPHOR"`
	MultiDisplay          bool    `vts:"MULTI_DISPLAY" json:"MULTI_DISPLAY"`
	OcToneMap             bool    `vts:"OC_TONEMAP" json:"OC_TONEMAP"`
	OcLightSamples        int     `vts:"OC_LIGHT_SAMPLES" json:"OC_LIGHT_SAMPLES"`
	OcDownSampleFactor    int     `vts:"OC_DOWNSAMPLE_FACTOR" json:"OC_DOWNSAMPLE_FACTOR"`
	BgmVolume             int     `vts:"BGM_VOLUME" json:"BGM_VOLUME"`
	VoiceVolume           int     `vts:"VOICE_VOLUME" json:"VOICE_VOLUME"`
	SkeletonFingers       bool    `vts:"SKELETON_FINGERS" json:"SKELETON_FINGERS"`
	TestQuickSave         bool    `vts:"TEST_QUICKSAVE" json:"TEST_QUICKSAVE"`
	CloudDiagnostics      bool    `vts:"CLOUD_DIAGNOSTICS" json:"CLOUD_DIAGNOSTICS"`
	PersistentPlayArea    bool    `vts:"PERSISTENT_PLAYAREA" json:"PERSISTENT_PLAYAREA"`
	LoadedVoicesVersion   string  `vts:"loadedVoicesVersion" json:"loadedVoicesVersion"`
	EulaAgreed            bool    `vts:"EULA_AGREED" json:"EULA_AGREED"`
	EulaNotif             bool    `vts:"EULA_NOTIF" json:"EULA_NOTIF"`
	CloudDiagnosticsNotif bool    `vts:"CLOUD_DIAGNOSTICS_NOTIF" json:"CLOUD_DIAGNOSTICS_NOTIF"`
	MpMyTailArt           string  `vts:"MP_MY_TAIL_ART" json:"MP_MY_TAIL_ART"`
	MpShowLiveries        string  `vts:"MP_SHOW_LIVERIES" json:"MP_SHOW_LIVERIES"`
	OnlineConduct         bool    `vts:"onlineConduct" json:"onlineConduct"`
	OnlineConductTime     string  `vts:"onlineConduct_time" json:"onlineConduct_time"`
	EquipSymmetry         bool    `vts:"equipSymmetry" json:"equipSymmetry"`
	NvgIpdOffSet          float32 `vts:"NVGIPDOffset" json:"NVGIPDOffset"`
	NvgBrightnessMul      float32 `vts:"NVGBrightnessMul" json:"NVGBrightnessMul"`
	NvgVerticalOffset     float32 `vts:"NVGVerticalOffset" json:"NVGVerticalOffset"`
	HostUnitIcons         bool    `vts:"host_unitIcons" json:"host_unitIcons"`
	HostLateJoins         bool    `vts:"host_lateJoins" json:"host_lateJoins"`
	HostCustomLiveries    bool    `vts:"host_customLiveries" json:"host_customLiveries"`
	HostSpectatorOption   string  `vts:"host_spectatorOption" json:"host_spectatorOption"`
	TwistRudderDeadzone   int     `vts:"TWIST_RUDDER_DEADZONE" json:"TWIST_RUDDER_DEADZONE"`
	Shadows               bool    `vts:"SHADOWS" json:"SHADOWS"`
	JetBorneAd            bool    `vts:"JETBORNE_AD" json:"JETBORNE_AD"`
	PlayArea              struct {
		PlayAreaPosition [3]float64 `vts:"playAreaPosition,tuple" json:"playAreaPosition"`
		PlayAreaRotation [3]float64 `vts:"playAreaRotation,tuple" json:"playAreaRotation"`
	} `vts:"PLAYAREA" json:"PLAYAREA"`
}

type BanFile struct {
	Time time.Time
	Bans []Ban
}

type Ban struct {
	SteamID uint64
	Reason  string
	BanDays int
}

type RecentPlayersFile struct {
	Users []User
}

type User struct {
	ID            uint64
	SteamName     string
	PilotName     string
	Time          time.Time
	TimeFirstSeen time.Time
}
