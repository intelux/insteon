package plm

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"
)

// DeviceInfo contains information about a device.
type DeviceInfo struct {
	X10HouseCode  byte
	X10Unit       byte
	RampRate      time.Duration
	OnLevel       float64
	LEDBrightness float64
}

func deviceInfoFromUserData(userData UserData) DeviceInfo {
	return DeviceInfo{
		X10HouseCode:  userData[4],
		X10Unit:       userData[5],
		RampRate:      byteToRampRate(userData[6]),
		OnLevel:       byteToOnLevel(userData[7]),
		LEDBrightness: byteToLEDBrightness(userData[8]),
	}
}

func deviceInfoToUserData(deviceInfo DeviceInfo) UserData {
	userData := UserData{}
	userData[4] = deviceInfo.X10HouseCode
	userData[5] = deviceInfo.X10Unit
	userData[6] = rampRateToByte(deviceInfo.RampRate)
	userData[7] = onLevelToByte(deviceInfo.OnLevel)
	userData[8] = ledBrightnessToByte(deviceInfo.LEDBrightness)

	return userData
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

// AllLinkRecordFlags represents an all-link record flags.
type AllLinkRecordFlags byte

// CommandBytes represent a pair of command bytes.
type CommandBytes [2]byte

func (b CommandBytes) String() string {
	switch b[0] {
	case 0x01:
		return fmt.Sprintf("assigned to all-link group %02d", b[1])
	case 0x02:
		return fmt.Sprintf("deleted from all-link group %02d", b[1])
	case 0x03:
		switch b[1] {
		case 0x00:
			return "product data requested"
		case 0x01:
			return "fx username requested"
		}
	case 0x09:
		return fmt.Sprintf("entered all-link linking for group %02d", b[1])
	case 0x0a:
		return fmt.Sprintf("entered all-link unlinking for group %02d", b[1])
	case 0x0d:
		switch b[1] {
		case 0x00:
			return "Insteon engine version requested"
		}
	case 0x0f:
		return "ping"
	case 0x10:
		return "id request"
	case 0x11:
		return fmt.Sprintf("turn on (level %.02f%%)", byteToOnLevel(b[1])*100)
	case 0x12:
		return fmt.Sprintf("turn on instantly (level %.02f%%)", byteToOnLevel(b[1])*100)
	case 0x13:
		return fmt.Sprintf("turn off (level %.02f%%)", byteToOnLevel(b[1])*100)
	case 0x14:
		return "turn off instantly"
	case 0x15:
		return "brighten one step"
	case 0x16:
		return "dim one step"
	case 0x17:
		switch b[1] {
		case 0x00:
			return "change started (dimming)"
		case 0x01:
			return "changed started (brightening)"
		}
	case 0x18:
		return "change stopped"
	case 0x19:
		switch b[1] {
		case 0x00:
			return "light status request (on level)"
		case 0x01:
			return "light status request (led info)"
		}
	case 0x1f:
		return "get operating flags"
	}

	return fmt.Sprintf("unknown command: %s", hex.EncodeToString(b[:]))
}

// UserData represent user data.
type UserData [14]byte

// LinkData represent link data.
type LinkData [3]byte

// AllLinkRecord represents a all-link record.
type AllLinkRecord struct {
	Flags    AllLinkRecordFlags
	Group    Group
	Identity Identity
	LinkData LinkData
}

// Mode returns the mode of an all-link record.
func (r AllLinkRecord) Mode() AllLinkMode {
	if r.Flags&0x40 > 0 {
		return ModeResponder
	}

	return ModeController
}

// AllLinkRecordList represents a list of all-link records.
type AllLinkRecordList []AllLinkRecord

// Len returns the length of the list.
func (l AllLinkRecordList) Len() int {
	return len(l)
}

// Less returns whether the element at i should appear before the element at j.
func (l AllLinkRecordList) Less(i, j int) bool {
	order := bytes.Compare(l[i].Identity[:], l[j].Identity[:])

	if order != 0 {
		return order < 0
	}

	if l[i].Group != l[j].Group {
		return l[i].Group < l[j].Group
	}

	return l[i].Mode() < l[j].Mode()
}

// Swap swaps two elements.
func (l AllLinkRecordList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// AllLinkMode represents an all-link mode.
type AllLinkMode byte

const (
	// ModeResponder represents a responder.
	ModeResponder AllLinkMode = 0x00
	// ModeController represents a controller.
	ModeController AllLinkMode = 0x01
	// ModeAuto represents auto-selection of the mode.
	ModeAuto AllLinkMode = 0x03
	// ModeDelete represents a deletion of an all-link record.
	ModeDelete AllLinkMode = 0xff
)

func (m AllLinkMode) String() string {
	switch m {
	case ModeResponder:
		return "responder"
	case ModeController:
		return "controller"
	case ModeAuto:
		return "auto"
	case ModeDelete:
		return "delete"
	default:
		panic(fmt.Errorf("unknown all-link mode %d", m))
	}
}

var (
	// CommandBytesBeep is used to make a device beep.
	CommandBytesBeep = CommandBytes{0x30, 0x00}

	// CommandBytesGetDeviceInfo is used to get the device information.
	CommandBytesGetDeviceInfo = CommandBytes{0x2e, 0x00}

	// CommandBytesStatusRequest is used to get the device status.
	CommandBytesStatusRequest = CommandBytes{0x19, 0x00}

	// CommandBytesSetDeviceInfo is used to set the device information.
	CommandBytesSetDeviceInfo = CommandBytes{0x2e, 0x00}
)
