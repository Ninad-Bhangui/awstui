package aws

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

func ParseProfilesFromConfig(configFilePath string) ([]string, error) {
	var profiles []string
	if configFilePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.New("Failed to get user home directory")
		}
		configFilePath = filepath.Join(homeDir, ".aws", "config")
	}

	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return profiles, nil
	}
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		return nil, errors.New("Could not load file")
	}
	sections := cfg.Sections()
	for _, section := range sections {
		if strings.HasPrefix(section.Name(), "profile ") {
			profile := strings.TrimPrefix(section.Name(), "profile ")
			profiles = append(profiles, profile)
		}
	}

	return profiles, nil

}
