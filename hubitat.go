package insteon

// HubitatConfiguration contains the Hubitat configuration.
type HubitatConfiguration struct {
	HubURL string `yaml:"hub_url"`
}

// HubitatEvent represents a Hubitat event.
type HubitatEvent struct {
	Alias string     `json:"id"`
	State LightState `json:"state"`
}
