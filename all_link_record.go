package insteon

import (
	"bytes"
	"fmt"
)

// AllLinkRecordFlags represents an all-link record flags.
type AllLinkRecordFlags byte

// AllLinkRecord represents a all-link record.
type AllLinkRecord struct {
	Flags    AllLinkRecordFlags `json:"flags"`
	Group    Group              `json:"group"`
	ID       ID                 `json:"id"`
	LinkData []byte             `json:"link_data"`
}

// UnmarshalBinary -
func (r *AllLinkRecord) UnmarshalBinary(b []byte) error {
	if len(b) != 8 {
		return fmt.Errorf("expected a buffer of size 8 but got one of size %d", len(b))
	}

	r.Flags = AllLinkRecordFlags(b[0])
	r.Group = Group(b[1])
	copy(r.ID[:], b[2:5])
	r.LinkData = make([]byte, 3)
	copy(r.LinkData, b[5:8])

	return nil
}

// MarshalBinary -
func (r *AllLinkRecord) MarshalBinary() ([]byte, error) {
	return []byte{
		byte(r.Flags),
		byte(r.Group),
		r.ID[0],
		r.ID[1],
		r.ID[2],
		r.LinkData[0],
		r.LinkData[1],
		r.LinkData[2],
	}, nil
}

// Mode returns the mode of an all-link record.
func (r AllLinkRecord) Mode() AllLinkMode {
	if r.Flags&0x40 > 0 {
		return ModeResponder
	}

	return ModeController
}

// AllLinkRecordSlice is a slice of all link records.
type AllLinkRecordSlice []AllLinkRecord

// Len returns the length of the list.
func (l AllLinkRecordSlice) Len() int {
	return len(l)
}

// Less returns whether the element at i should appear before the element at j.
func (l AllLinkRecordSlice) Less(i, j int) bool {
	order := bytes.Compare(l[i].ID[:], l[j].ID[:])

	if order != 0 {
		return order < 0
	}

	if l[i].Group != l[j].Group {
		return l[i].Group < l[j].Group
	}

	return l[i].Mode() < l[j].Mode()
}

// Swap swaps two elements.
func (l AllLinkRecordSlice) Swap(i, j int) {
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
