package plm

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"time"
)

// Identity is an Insteon identity.
type Identity [3]byte

func (i Identity) String() string {
	return hex.EncodeToString(i[:])
}

// AsGroup returns the identity as a group. Result is only valid if the
// identity is the target of a broadcast message.
func (i Identity) AsGroup() Group {
	return Group(i[2])
}

// ParseIdentity parses an identity.
func ParseIdentity(s string) (Identity, error) {
	var identity Identity

	b, err := hex.DecodeString(s)

	if err != nil {
		return identity, err
	}

	if len(b) != 3 {
		return identity, fmt.Errorf("invalid identity (%s)", s)
	}

	copy(identity[:], b)

	return identity, nil
}

// MainCategory represents a main category.
type MainCategory uint8

var (
	generalizedControllers  MainCategory = 0x00
	dimmableLightingControl MainCategory = 0x01
	switchedLightingControl MainCategory = 0x02
	networkBridges          MainCategory = 0x03
	irrigationControl       MainCategory = 0x04
	climateControlHeating   MainCategory = 0x05
	poolAndSpaControl       MainCategory = 0x06
	sensorsAndActuators     MainCategory = 0x07
	homeEntertainement      MainCategory = 0x08
	energyManagement        MainCategory = 0x09
	builtInApplianceControl MainCategory = 0x0A
	plumbing                MainCategory = 0x0B
	communication           MainCategory = 0x0C
	computerControl         MainCategory = 0x0D
	windowCoverings         MainCategory = 0x0E
	accessControl           MainCategory = 0x0F
	securityHealthSafety    MainCategory = 0x10
	surveillance            MainCategory = 0x11
	automotive              MainCategory = 0x12
	petCare                 MainCategory = 0x13
	toys                    MainCategory = 0x14
	timekeeping             MainCategory = 0x15
	holiday                 MainCategory = 0x16
	unassigned              MainCategory = 0xFF

	// networkBridges subcategories.
	powerlincSerial               SubCategory = 0x01
	powerlincUsb                  SubCategory = 0x02
	iconPowerlincSerial           SubCategory = 0x03
	iconPowerlincUsb              SubCategory = 0x04
	smartlabsPowerLineModemSerial SubCategory = 0x05
	powerlincDualBandSerial       SubCategory = 0x11
	powerlincDualBandUsb          SubCategory = 0x15
)

// SubCategory represents a main category.
type SubCategory uint8

// Category represents a category.
type Category struct {
	mainCategory MainCategory
	subCategory  SubCategory
}

func (c Category) String() string {
	switch c.mainCategory {
	case generalizedControllers:
		return "Generalized Controllers"
	case dimmableLightingControl:
		return "Dimmable Lighting Control"
	case switchedLightingControl:
		return "Switched Lighting Control"
	case networkBridges:
		switch c.subCategory {
		case powerlincSerial:
			return "PowerLinc Serial [2414S]"
		case powerlincUsb:
			return "PowerLinc USB [2414U]"
		case iconPowerlincSerial:
			return "Icon PowerLinc Serial [2814 S]"
		case iconPowerlincUsb:
			return "Icon PowerLinc USB [2814U] "
		case smartlabsPowerLineModemSerial:
			return "Smartlabs Power Line Modem Serial [2412S]"
		case powerlincDualBandSerial:
			return "PowerLinc Dual Band Serial [2413S]"
		case powerlincDualBandUsb:
			return "PowerLinc Dual Band USB [2413U]"
		}

		return "Network Bridges"
	case irrigationControl:
		return "Irrigation Control"
	case climateControlHeating:
		return "Climate Control"
	case poolAndSpaControl:
		return "Pool and Spa Control"
	case sensorsAndActuators:
		return "Sensors and Actuators"
	case homeEntertainement:
		return "Home Entertainment"
	case energyManagement:
		return "Energy Management"
	case builtInApplianceControl:
		return "Built-In Appliance Control"
	case plumbing:
		return "Plumbing"
	case communication:
		return "Communication"
	case computerControl:
		return "Computer Control"
	case windowCoverings:
		return "Window Coverings"
	case accessControl:
		return "Access Control"
	case securityHealthSafety:
		return "Security Health Safety"
	case surveillance:
		return "Surveillance"
	case automotive:
		return "Automotive"
	case petCare:
		return "Pet Care"
	case toys:
		return "Toys"
	case timekeeping:
		return "Timekeeping"
	case holiday:
		return "Holiday"
	case unassigned:
		return "Unassigned"
	}

	return "Unknown category"
}

// IMInfo contains information about the PLM.
type IMInfo struct {
	Identity        Identity
	Category        Category
	FirmwareVersion uint8
}

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

func clampLevel(v float64) float64 {
	return math.Max(0, math.Min(1, v))
}

func byteToOnLevel(b byte) float64 {
	return float64(b) / 0xff
}

func onLevelToByte(level float64) byte {
	return byte(clampLevel(level) * 0xff)
}

func byteToLEDBrightness(b byte) float64 {
	return float64(b&0x7f) / 0x7f
}

func ledBrightnessToByte(level float64) byte {
	return byte(clampLevel(level) * 0x7f)
}

// LightOnOff represents a light on/off state.
type LightOnOff bool

const (
	// LightOn indicates an on light.
	LightOn LightOnOff = true
	// LightOff indicates an off light.
	LightOff LightOnOff = false
)

func (s LightOnOff) String() string {
	if s {
		return "on"
	}

	return "off"
}

// LightStateChange represents a light state change.
type LightStateChange int

const (
	// ChangeNormal indicates that the light state must change as if a single
	// press had been done.
	ChangeNormal LightStateChange = iota
	// ChangeInstant indicates that the light state must change instantly, as
	// if a quick double-press had been done.
	ChangeInstant
	// ChangeStep indicates that the light state must change for one step up or
	// down.
	ChangeStep
	// ChangeStart indicates that the light state must start changing until a
	// ChangeStop change is set.
	ChangeStart
	// ChangeStop stops a light change started with ChangeStart.
	ChangeStop
)

// LightState represents a light state.
type LightState struct {
	OnOff  LightOnOff
	Change LightStateChange
	Level  float64
}

func (s LightState) commandBytes() CommandBytes {
	levelByte := onLevelToByte(s.Level)

	if s.OnOff == LightOn {
		switch s.Change {
		case ChangeInstant:
			return CommandBytes([2]byte{0x12, levelByte})
		case ChangeStep:
			return CommandBytes([2]byte{0x15, 0})
		case ChangeStart:
			return CommandBytes([2]byte{0x17, 0x01})
		case ChangeStop:
			return CommandBytes([2]byte{0x18, 0x00})
		}

		return CommandBytes([2]byte{0x11, levelByte})
	}

	switch s.Change {
	case ChangeInstant:
		return CommandBytes([2]byte{0x14, levelByte})
	case ChangeStep:
		return CommandBytes([2]byte{0x16, 0})
	case ChangeStart:
		return CommandBytes([2]byte{0x17, 0x00})
	case ChangeStop:
		return CommandBytes([2]byte{0x18, 0x00})
	}

	return CommandBytes([2]byte{0x13, levelByte})
}

// CommandBytesToLightState get light state from command bytes.
func CommandBytesToLightState(bytes CommandBytes) (state *LightState) {
	switch bytes[0] {
	case 0x11:
		state = &LightState{
			OnOff:  LightOn,
			Change: ChangeNormal,
			Level:  byteToOnLevel(bytes[1]),
		}
	case 0x12:
		state = &LightState{
			OnOff:  LightOn,
			Change: ChangeInstant,
			Level:  byteToOnLevel(bytes[1]),
		}
	case 0x13:
		state = &LightState{
			OnOff:  LightOff,
			Change: ChangeNormal,
			Level:  byteToOnLevel(bytes[1]),
		}
	case 0x14:
		state = &LightState{
			OnOff:  LightOff,
			Change: ChangeInstant,
			Level:  byteToOnLevel(bytes[1]),
		}
	case 0x15:
		state = &LightState{
			OnOff:  LightOn,
			Change: ChangeStep,
			Level:  0,
		}
	case 0x16:
		state = &LightState{
			OnOff:  LightOff,
			Change: ChangeStep,
			Level:  0,
		}
	case 0x17:
		if bytes[1] == 0x00 {
			state = &LightState{
				OnOff:  LightOff,
				Change: ChangeStart,
				Level:  0,
			}
		} else {
			state = &LightState{
				OnOff:  LightOn,
				Change: ChangeStart,
				Level:  0,
			}
		}
	case 0x18:
		state = &LightState{
			OnOff:  LightOn,
			Change: ChangeStop,
			Level:  0,
		}
	}
	return
}

// MessageFlags represents the message flags.
type MessageFlags byte

const (
	// MessageFlagExtended indicates extended messages.
	MessageFlagExtended MessageFlags = 0x10
	// MessageFlagAck indicates an acquitement message.
	MessageFlagAck MessageFlags = 0x20
	// MessageFlagAllLink indicates an all-link message.
	MessageFlagAllLink MessageFlags = 0x40
	// MessageFlagBroadcast indicates a broadcast message.
	MessageFlagBroadcast MessageFlags = 0x80
)

// AllLinkRecordFlags represents an all-link record flags.
type AllLinkRecordFlags byte

// Group represents a group.
type Group byte

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

	// CommandBytesSetDeviceInfo is used to set the device information.
	CommandBytesSetDeviceInfo = CommandBytes{0x2e, 0x00}
)

func checksum(commandBytes CommandBytes, userData UserData) byte {
	var checksum byte

	for _, b := range commandBytes {
		checksum += b
	}
	for _, b := range userData {
		checksum += b
	}

	return ((0xff ^ checksum) + 1) & 0xff
}
