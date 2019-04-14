package insteon

import "fmt"

// IMInfo contains information about a PowerLine Modem.
type IMInfo struct {
	ID              ID       `json:"id"`
	Category        Category `json:"category"`
	FirmwareVersion uint8    `json:"firmware_version"`
}

// UnmarshalBinary -
func (i *IMInfo) UnmarshalBinary(b []byte) error {
	if len(b) != 6 {
		return fmt.Errorf("expected 6 bytes but got %d", len(b))
	}

	copy(i.ID[:], b[:3])
	i.Category.UnmarshalBinary(b[3:5])
	i.FirmwareVersion = b[5]

	return nil
}

// MarshalBinary -
func (i IMInfo) MarshalBinary() ([]byte, error) {
	b := make([]byte, 6)
	copy(b[:3], i.ID[:])
	cb, _ := i.Category.MarshalBinary()
	copy(b[3:5], cb)
	b[5] = i.FirmwareVersion

	return b, nil
}
