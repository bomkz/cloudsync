package pilots

import (
	"time"

	"github.com/bomkz/cloudsync/pilots/vtscfg"
)

type PilotsFile struct {
	Pilots Pilots `vts:"PILOTS"`
}

type Pilots struct {
	PilotSaves []PilotSave `vts:"PILOTSAVE"`
}

type PilotSave struct {
	PilotName       string     `vts:"pilotName"`
	LastVehicleUsed string     `vts:"lastVehicleUsed"`
	TotalFlightTime float64    `vts:"totalFlightTime"`
	SkinColor       [3]float64 `vts:"skinColor"`
	SuitColor       [3]float64 `vts:"suitColor"`
	VestColor       [3]float64 `vts:"vestColor"`
	GSuitColor      [3]float64 `vts:"gSuitColor"`
	StrapsColor     [3]float64 `vts:"strapsColor"`
	Vehicles        []Vehicle  `vts:"VEHICLE"`
}

type Vehicle struct {
	VehicleName      string        `vts:"vehicleName"`
	SeatHeight       float64       `vts:"seatHeight"`
	JoystickPosition [3]float64    `vts:"joystickPosition"`
	ThrottlePosition [3]float64    `vts:"throttlePosition"`
	AltitudeMode     string        `vts:"altitudeMode"`
	DistanceMode     string        `vts:"distanceMode"`
	AirspeedMode     string        `vts:"airspeedMode"`
	LastLiveryID     string        `vts:"lastLiveryId"`
	VData            vtscfg.Node   `vts:"VDATA"`
	Campaigns        []Campaign    `vts:"CAMPAIGN"`
	SavedLoadouts    SavedLoadouts `vts:"SavedLoadouts"`
}

type Campaign struct {
	CampaignName         string         `vts:"campaignName"`
	CampaignID           string         `vts:"campaignID"`
	VehicleName          string         `vts:"vehicleName"`
	AvailableWeapons     string         `vts:"availableWeapons"`
	CurrentFuel          float64        `vts:"currentFuel"`
	AvailableScenarios   string         `vts:"availableScenarios"`
	LastScenarioIndex    int            `vts:"lastScenarioIdx"`
	LastScenarioTraining bool           `vts:"lastScenarioWasTraining"`
	CurrentWeapons       CurrentWeapons `vts:"currentWeapons"`
}

type CurrentWeapons struct {
	Slots []WeaponSlot `vts:"weapon"`
}

type WeaponSlot struct {
	Index  int    `vts:"idx"`
	Weapon string `vts:"weapon"`
}

type SavedLoadouts struct {
	Loadouts []SavedLoadout `vts:"SavedLoadout"`
}

type SavedLoadout struct {
	Name     string  `vts:"name"`
	NormFuel float64 `vts:"normFuel"`
	Eq0      string  `vts:"eq0"`
	Eq1      string  `vts:"eq1"`
	Eq2      string  `vts:"eq2"`
	Eq3      string  `vts:"eq3"`
	Eq4      string  `vts:"eq4"`
	Eq5      string  `vts:"eq5"`
	Eq6      string  `vts:"eq6"`
	Eq7      string  `vts:"eq7"`
	Eq8      string  `vts:"eq8"`
	Eq9      string  `vts:"eq9"`
	Eq10     string  `vts:"eq10"`
	Eq11     string  `vts:"eq11"`
	Eq12     string  `vts:"eq12"`
	Eq13     string  `vts:"eq13"`
	Eq14     string  `vts:"eq14"`
	Eq15     string  `vts:"eq15"`
}

type GameSettings struct {
	Settings GameSettingsWrapper `vts:"GAMESETTINGS"`
}
type GameSettingsWrapper struct {
	RadioMusicPath        string  `vts:"RADIO_MUSIC_PATH"`
	WingmanVoices         string  `vts:"WINGMAN_VOICES"`
	ToolTips              bool    `vts:"TOOLTIPS"`
	UnitIcons             bool    `vts:"UNIT_ICONS"`
	BodyPhysics           bool    `vts:"BODY_PHYSICS"`
	TreeCollisions        bool    `vts:"TREE_COLLISIONS"`
	HookPhysics           bool    `vts:"HOOK_PHYSICS"`
	PersistentSCam        bool    `vts:"PERSISTENT_S_CAM"`
	ShowBobbleHead        bool    `vts:"SHOW_BOBBLEHEAD"`
	ThumbstickMode        bool    `vts:"THUMBSTICK_MODE"`
	ThumbstickDeadZone    int     `vts:"THUMBSTICK_DEADZONE"`
	ThumbRudder           bool    `vts:"THUMB_RUDDER"`
	HardwareRudder        bool    `vts:"HARDWARE_RUDDER"`
	TapToggleGrip         bool    `vts:"TAP_TOGGLE_GRIP"`
	ControlHaptics        int     `vts:"CONTROL_HAPTICS"`
	OverallHaptics        int     `vts:"OVERALL_HAPTICS"`
	Msaa                  int     `vts:"MSAA"`
	HideHelmet            bool    `vts:"HIDE_HELMET"`
	FullScreenNVG         bool    `vts:"FULLSCREEN_NVG"`
	NvgPhosphor           bool    `vts:"NVG_PHOSPHOR"`
	MultiDisplay          bool    `vts:"MULTI_DISPLAY"`
	OcToneMap             bool    `vts:"OC_TONEMAP"`
	OcLightSamples        int     `vts:"OC_LIGHT_SAMPLES"`
	OcDownSampleFactor    int     `vts:"OC_DOWNSAMPLE_FACTOR"`
	BgmVolume             int     `vts:"BGM_VOLUME"`
	VoiceVolume           int     `vts:"VOICE_VOLUME"`
	SkeletonFingers       bool    `vts:"SKELETON_FINGERS"`
	TestQuickSave         bool    `vts:"TEST_QUICKSAVE"`
	CloudDiagnostics      bool    `vts:"CLOUD_DIAGNOSTICS"`
	PersistentPlayArea    bool    `vts:"PERSISTENT_PLAYAREA"`
	LoadedVoicesVersion   string  `vts:"loadedVoicesVersion"`
	EulaAgreed            bool    `vts:"EULA_AGREED"`
	EulaNotif             bool    `vts:"EULA_NOTIF"`
	CloudDiagnosticsNotif bool    `vts:"CLOUD_DIAGNOSTICS_NOTIF"`
	MpMyTailArt           string  `vts:"MP_MY_TAIL_ART"`
	MpShowLiveries        string  `vts:"MP_SHOW_LIVERIES"`
	OnlineConduct         bool    `vts:"onlineConduct"`
	OnlineConductTime     string  `vts:"onlineConduct_time"`
	EquipSymmetry         bool    `vts:"equipSymmetry"`
	NvgIpdOffSet          float32 `vts:"NVGIPDOffset"`
	NvgBrightnessMul      float32 `vts:"NVGBrightnessMul"`
	NvgVerticalOffset     float32 `vts:"NVGVerticalOffset"`
	HostUnitIcons         bool    `vts:"host_unitIcons"`
	HostLateJoins         bool    `vts:"host_lateJoins"`
	HostCustomLiveries    bool    `vts:"host_customLiveries"`
	HostSpectatorOption   string  `vts:"host_spectatorOption"`
	TwistRudderDeadzone   int     `vts:"TWIST_RUDDER_DEADZONE"`
	Shadows               bool    `vts:"SHADOWS"`
	JetBorneAd            bool    `vts:"JETBORNE_AD"`
	PlayArea              struct {
		PlayAreaPosition [3]float64 `vts:"playAreaPosition,tuple"`
		PlayAreaRotation [3]float64 `vts:"playAreaRotation,tuple"`
	} `vts:"PLAYAREA"`
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
