package insteon

// DeviceEvent represents a DeviceEvent.
type DeviceEvent struct {
	Identity ID               `json:"id"`
	OnOff    LightOnOff       `json:"onoff"`
	Change   LightStateChange `json:"change,omitempty"`
}
