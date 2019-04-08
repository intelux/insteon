package insteon

// IMInfo contains information about a PowerLine Modem.
type IMInfo struct {
	ID              ID       `json:"id"`
	Category        Category `json:"category"`
	FirmwareVersion uint8    `json:"firmware_version"`
}
