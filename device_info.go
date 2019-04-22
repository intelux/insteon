package insteon

import (
	"fmt"
	"time"
)

// DeviceInfo contains information about a device.
type DeviceInfo struct {
	X10Address    *[2]byte       `json:"x10_address,omitempty"`
	RampRate      *time.Duration `json:"ramp_rate,omitempty"`
	OnLevel       *float64       `json:"on_level,omitempty"`
	LEDBrightness *float64       `json:"led_brightness,omitempty"`
}

// UnmarshalBinary -
func (i *DeviceInfo) UnmarshalBinary(b []byte) error {
	if len(b) != 14 {
		return fmt.Errorf("expected 14 bytes but got %d", len(b))
	}

	x10Address := [2]byte{b[4], b[5]}
	i.X10Address = &x10Address
	rampRate := byteToRampRate(b[6])
	i.RampRate = &rampRate
	onLevel := byteToOnLevel(b[7])
	i.OnLevel = &onLevel
	ledBrightness := byteToLEDBrightness(b[8])
	i.LEDBrightness = &ledBrightness

	return nil
}

// MarshalBinary -
func (i DeviceInfo) MarshalBinary() ([]byte, error) {
	result := make([]byte, 14)
	result[4] = (*i.X10Address)[0]
	result[5] = (*i.X10Address)[1]
	result[6] = rampRateToByte(*i.RampRate)
	result[7] = onLevelToByte(*i.OnLevel)
	result[8] = ledBrightnessToByte(*i.LEDBrightness)

	return result, nil
}

var rampRates = []struct {
	Duration time.Duration
	Byte     byte
}{
	{time.Millisecond * 100, 0x1f},
	{time.Millisecond * 200, 0x1e},
	{time.Millisecond * 300, 0x1d},
	{time.Millisecond * 500, 0x1c},
	{time.Second * 2, 0x1b},
	{time.Millisecond * 4500, 0x1a},
	{time.Millisecond * 6500, 0x19},
	{time.Millisecond * 8500, 0x18},
	{time.Second * 19, 0x17},
	{time.Millisecond * 21500, 0x16},
	{time.Millisecond * 23500, 0x15},
	{time.Second * 26, 0x14},
	{time.Second * 28, 0x13},
	{time.Second * 30, 0x12},
	{time.Second * 32, 0x11},
	{time.Second * 34, 0x10},
	{time.Millisecond * 38500, 0x0f},
	{time.Second * 43, 0x0e},
	{time.Second * 47, 0x0d},
	{time.Second * 60, 0x0c},
	{time.Second * 90, 0x0b},
	{time.Second * 120, 0x0a},
	{time.Second * 150, 0x09},
	{time.Second * 180, 0x08},
	{time.Second * 210, 0x07},
	{time.Second * 240, 0x06},
	{time.Second * 270, 0x05},
	{time.Second * 300, 0x04},
	{time.Second * 360, 0x03},
	{time.Second * 420, 0x02},
	{time.Second * 480, 0x01},
}

func byteToRampRate(b byte) time.Duration {
	var value = rampRates[0].Duration

	for _, rampRate := range rampRates {
		if b > rampRate.Byte {
			break
		}

		value = rampRate.Duration
	}

	return value
}

func rampRateToByte(duration time.Duration) byte {
	var value = rampRates[0].Byte

	for _, rampRate := range rampRates {
		if duration < rampRate.Duration {
			break
		}

		value = rampRate.Byte
	}

	return value
}

func byteToLEDBrightness(b byte) float64 {
	return float64(b&0x7f) / 0x7f
}

func ledBrightnessToByte(level float64) byte {
	return byte(clampLevel(level) * 0x7f)
}
