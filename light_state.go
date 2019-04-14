package insteon

import (
	"fmt"
	"math"
)

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

func (s LightState) asCommandBytes() [2]byte {
	data, _ := s.MarshalBinary()

	return [2]byte{data[0], data[1]}
}

// MarshalBinary -
func (s LightState) MarshalBinary() ([]byte, error) {
	levelByte := onLevelToByte(s.Level)

	if s.OnOff == LightOn {
		switch s.Change {
		case ChangeInstant:
			return []byte{0x12, levelByte}, nil
		case ChangeStep:
			return []byte{0x15, 0}, nil
		case ChangeStart:
			return []byte{0x17, 0x01}, nil
		case ChangeStop:
			return []byte{0x18, 0x00}, nil
		}

		return []byte{0x11, levelByte}, nil
	}

	switch s.Change {
	case ChangeInstant:
		return []byte{0x14, levelByte}, nil
	case ChangeStep:
		return []byte{0x16, 0}, nil
	case ChangeStart:
		return []byte{0x17, 0x00}, nil
	case ChangeStop:
		return []byte{0x18, 0x00}, nil
	}

	return []byte{0x13, levelByte}, nil
}

// UnmarshalBinary -
func (s *LightState) UnmarshalBinary(b []byte) error {
	if len(b) != 2 {
		return fmt.Errorf("expected 2 bytes, not %d", len(b))
	}

	switch b[0] {
	case 0x11:
		*s = LightState{
			OnOff:  LightOn,
			Change: ChangeNormal,
			Level:  byteToOnLevel(b[1]),
		}
	case 0x12:
		*s = LightState{
			OnOff:  LightOn,
			Change: ChangeInstant,
			Level:  byteToOnLevel(b[1]),
		}
	case 0x13:
		*s = LightState{
			OnOff:  LightOff,
			Change: ChangeNormal,
			Level:  byteToOnLevel(b[1]),
		}
	case 0x14:
		*s = LightState{
			OnOff:  LightOff,
			Change: ChangeInstant,
			Level:  byteToOnLevel(b[1]),
		}
	case 0x15:
		*s = LightState{
			OnOff:  LightOn,
			Change: ChangeStep,
			Level:  0,
		}
	case 0x16:
		*s = LightState{
			OnOff:  LightOff,
			Change: ChangeStep,
			Level:  0,
		}
	case 0x17:
		if b[1] == 0x00 {
			*s = LightState{
				OnOff:  LightOff,
				Change: ChangeStart,
				Level:  0,
			}
		} else {
			*s = LightState{
				OnOff:  LightOn,
				Change: ChangeStart,
				Level:  0,
			}
		}
	case 0x18:
		*s = LightState{
			OnOff:  LightOn,
			Change: ChangeStop,
			Level:  0,
		}
	default:
		return fmt.Errorf("unexpected command code for light state: %02x%02x", b[0], b[1])
	}

	return nil
}

func byteToOnLevel(b byte) float64 {
	return float64(b) / 0xff
}

func onLevelToByte(level float64) byte {
	return byte(clampLevel(level) * 0xff)
}

func clampLevel(v float64) float64 {
	return math.Max(0, math.Min(1, v))
}
