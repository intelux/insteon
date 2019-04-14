package plm

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

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
