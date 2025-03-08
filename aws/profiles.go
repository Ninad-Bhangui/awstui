package aws

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
)

// Profile represents an AWS profile configuration
type Profile struct {
	Name      string
	Region    string
	IsDefault bool
	IsFromEnv bool
}

// GetProfiles returns a list of available AWS profiles from both config and credentials files
func GetProfiles() ([]Profile, error) {
	var profiles []Profile
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.New("failed to get user home directory")
	}

	// Get profiles from config file
	configProfiles, err := parseProfilesFromFile(filepath.Join(homeDir, ".aws", "config"), true)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	profiles = append(profiles, configProfiles...)

	// Get profiles from credentials file
	credProfiles, err := parseProfilesFromFile(filepath.Join(homeDir, ".aws", "credentials"), false)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	profiles = append(profiles, credProfiles...)

	// Deduplicate profiles and mark default
	profiles = deduplicateProfiles(profiles)

	// Check AWS_PROFILE environment variable
	if envProfile := os.Getenv("AWS_PROFILE"); envProfile != "" {
		for i := range profiles {
			if profiles[i].Name == envProfile {
				profiles[i].IsFromEnv = true
				break
			}
		}
	}

	return profiles, nil
}

func parseProfilesFromFile(filePath string, isConfig bool) ([]Profile, error) {
	var profiles []Profile

	_, err := os.Stat(filePath)
	if err != nil {
		return profiles, err
	}

	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, errors.New("could not load file: " + filePath)
	}

	sections := cfg.Sections()
	for _, section := range sections {
		name := section.Name()
		if name == "DEFAULT" || name == ini.DefaultSection {
			continue
		}

		if isConfig {
			// In config file, profiles are prefixed with "profile " except for 'default'
			if name == "default" {
				profiles = append(profiles, Profile{
					Name:      "default",
					Region:    section.Key("region").String(),
					IsDefault: true,
				})
			} else if strings.HasPrefix(name, "profile ") {
				profiles = append(profiles, Profile{
					Name:   strings.TrimPrefix(name, "profile "),
					Region: section.Key("region").String(),
				})
			}
		} else {
			// In credentials file, profile names are used directly
			profiles = append(profiles, Profile{
				Name:      name,
				IsDefault: name == "default",
			})
		}
	}

	return profiles, nil
}

func deduplicateProfiles(profiles []Profile) []Profile {
	seen := make(map[string]int)
	result := make([]Profile, 0)

	for _, p := range profiles {
		if idx, exists := seen[p.Name]; exists {
			// If profile exists and has a region, update the existing one
			if p.Region != "" {
				result[idx].Region = p.Region
			}
			// Preserve IsDefault flag
			result[idx].IsDefault = result[idx].IsDefault || p.IsDefault
		} else {
			seen[p.Name] = len(result)
			result = append(result, p)
		}
	}

	return result
}
