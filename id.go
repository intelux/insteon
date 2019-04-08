package insteon

import (
	"encoding/hex"
	"fmt"
)

// ID represents a device ID.
type ID [3]byte

func (i ID) String() string {
	return hex.EncodeToString(i[:])
}

// AsGroup returns the ID as a group. Result is only meaningful if the identity
// is the target of a broadcast message.
func (i ID) AsGroup() Group {
	return Group(i[2])
}

// UnmarshalText implements text unmarshalling.
func (i *ID) UnmarshalText(b []byte) error {
	data, err := hex.DecodeString(string(b))

	if err != nil {
		return fmt.Errorf("failed to hex-decode string: %s", err)
	}

	if len(data) != 3 {
		return fmt.Errorf("invalid size for identity: expected 3 but got %d byte(s)", len(data))
	}

	copy((*i)[:], data)

	return nil
}

// MarshalText implements text marshaling.
func (i ID) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// ParseID parses an ID.
func ParseID(s string) (id ID, err error) {
	err = id.UnmarshalText([]byte(s))

	return
}
