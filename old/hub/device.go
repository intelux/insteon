package hub

import "github.com/intelux/insteon/plm"

// Device represents a physical device.
type Device struct {
	ID          plm.Identity `json:"id" yaml:"id"`
	Name        string       `json:"name" yaml:"name"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
}
