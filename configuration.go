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

// ConfigurationDevice represents a device in the configuration.
type ConfigurationDevice struct {
	ID    ID     `yaml:"id"`
	Name  string `yaml:"name"`
	Group string `yaml:"group,omitempty"`
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
