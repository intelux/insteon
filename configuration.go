package insteon

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Configuration represents a configuration.
type Configuration struct {
	Devices []ConfigurationDevice `yaml:"devices"`
	Hubitat HubitatConfiguration  `yaml:"hubitat"`
}

// ErrNoSuchDevice is returned whenever a lookup on a given device alias
// failed.
type ErrNoSuchDevice struct {
	ID    ID
	Alias string
}

// Error returns the error string.
func (e ErrNoSuchDevice) Error() string {
	if e.Alias != "" {
		return fmt.Sprintf("no such device with alias: %s", e.Alias)
	}

	return fmt.Sprintf("no such device with id: %s", e.ID)
}

// GetDevice finds a device from its id.
func (c *Configuration) GetDevice(id ID) (*ConfigurationDevice, error) {
	for _, device := range c.Devices {
		if device.ID == id {
			return &device, nil
		}
	}

	return nil, ErrNoSuchDevice{ID: id}
}

// LookupDevice finds a device from its alias.
func (c *Configuration) LookupDevice(alias string) (*ConfigurationDevice, error) {
	for _, device := range c.Devices {
		if device.Alias == alias {
			return &device, nil
		}
	}

	return nil, ErrNoSuchDevice{Alias: alias}
}

// ConfigurationDevice represents a device in the configuration.
type ConfigurationDevice struct {
	ID              ID     `yaml:"id" json:"insteon_id"`
	Name            string `yaml:"name" json:"description"`
	Alias           string `yaml:"alias" json:"id"`
	Group           string `yaml:"group,omitempty" json:"-"`
	MirrorDeviceIDs []ID   `yaml:"mirror_devices" json:"-"`
	ControllerIDs   []ID   `yaml:"controllers" json:"-"`
}

// UnmarshalYAML -
func (d *ConfigurationDevice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type X ConfigurationDevice

	x := &X{}

	if err := unmarshal(x); err != nil {
		return err
	}

	if x.Name == "" {
		return fmt.Errorf("a name must be defined")
	}

	if x.Alias == "" {
		return fmt.Errorf("an alias must be defined")
	}

	*d = *(*ConfigurationDevice)(x)

	return nil
}

func getUserConfigPath() (string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, ".config", "ion", "config.yml"), nil
}

func getSystemConfigPath() string {
	return "/etc/ion/config.yml"
}

// LoadDefaultConfiguration loads the default configuration.
func LoadDefaultConfiguration() (*Configuration, error) {
	paths := make([]string, 0, 2)

	if usrPath, err := getUserConfigPath(); err == nil {
		paths = append(paths, usrPath)
	}

	paths = append(paths, getSystemConfigPath())

	for _, path := range paths {
		f, err := os.Open(path)

		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, err
		}

		defer f.Close()

		config, err := LoadConfiguration(f)

		if err != nil {
			return nil, fmt.Errorf("reading configuration at %s: %s", path, err)
		}

		return config, nil
	}

	return &Configuration{}, nil
}

// LoadConfiguration loads a configuration from a YAML stream.
func LoadConfiguration(r io.Reader) (*Configuration, error) {
	config := &Configuration{}

	if err := yaml.NewDecoder(r).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
