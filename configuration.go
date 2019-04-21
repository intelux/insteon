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
}

// LookupDevice finds a device from its alias or device ID.
//
// Devices that are not referenced in the configuration can still be looked-up
// by their device ID.
func (c *Configuration) LookupDevice(s string) (ID, error) {
	if s != "" {
		for _, device := range c.Devices {
			if device.Alias == s {
				return device.ID, nil
			}
		}
	}

	return ParseID(s)
}

// ConfigurationDevice represents a device in the configuration.
type ConfigurationDevice struct {
	ID    ID     `yaml:"id"`
	Name  string `yaml:"name"`
	Alias string `yaml:"alias,omitempty"`
	Group string `yaml:"group,omitempty"`
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
