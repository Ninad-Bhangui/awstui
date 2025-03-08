package aws

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProfiles(t *testing.T) {
	// Set up test directory structure
	tmpDir := t.TempDir()
	awsDir := filepath.Join(tmpDir, ".aws")
	os.MkdirAll(awsDir, 0755)

	// Create test config file
	configContent := `[default]
region = us-east-1

[profile dev]
region = us-west-2

[profile prod]
region = eu-west-1
`
	configPath := filepath.Join(awsDir, "config")
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)

	// Create test credentials file
	credentialsContent := `[default]
aws_access_key_id = default_key
aws_secret_access_key = default_secret

[dev]
aws_access_key_id = dev_key
aws_secret_access_key = dev_secret

[staging]
aws_access_key_id = staging_key
aws_secret_access_key = staging_secret
`
	credentialsPath := filepath.Join(awsDir, "credentials")
	err = os.WriteFile(credentialsPath, []byte(credentialsContent), 0644)
	assert.NoError(t, err)

	// Test without AWS_PROFILE set
	t.Run("Without AWS_PROFILE", func(t *testing.T) {
		// Temporarily set HOME to our test directory
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpDir)
		defer os.Setenv("HOME", originalHome)

		pm := NewProfileManager()
		err := pm.LoadProfiles()
		assert.NoError(t, err)
		profiles := pm.GetAllProfiles()

		// Should find 4 unique profiles: default, dev, prod, staging
		assert.Equal(t, 4, len(profiles))

		// Verify default profile
		var defaultProfile Profile
		for _, p := range profiles {
			if p.Name == "default" {
				defaultProfile = p
				break
			}
		}
		assert.True(t, defaultProfile.IsDefault)
		assert.Equal(t, "us-east-1", defaultProfile.Region)

		// Verify dev profile has region from config
		var devProfile Profile
		for _, p := range profiles {
			if p.Name == "dev" {
				devProfile = p
				break
			}
		}
		assert.Equal(t, "us-west-2", devProfile.Region)
	})

	// Test with AWS_PROFILE set
	t.Run("With AWS_PROFILE", func(t *testing.T) {
		// Temporarily set HOME and AWS_PROFILE
		originalHome := os.Getenv("HOME")
		originalProfile := os.Getenv("AWS_PROFILE")
		os.Setenv("HOME", tmpDir)
		os.Setenv("AWS_PROFILE", "dev")
		defer func() {
			os.Setenv("HOME", originalHome)
			os.Setenv("AWS_PROFILE", originalProfile)
		}()

		pm := NewProfileManager()
		err := pm.LoadProfiles()
		assert.NoError(t, err)
		profiles := pm.GetAllProfiles()

		// Verify dev profile is marked as from env
		var devProfile Profile
		for _, p := range profiles {
			if p.Name == "dev" {
				devProfile = p
				break
			}
		}
		assert.True(t, devProfile.IsFromEnv)
		assert.Equal(t, "us-west-2", devProfile.Region)
	})
}

func TestParseProfilesFromFile(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("Non-existent file", func(t *testing.T) {
		profiles, err := parseProfilesFromFile(filepath.Join(tmpDir, "nonexistent"), false)
		assert.Error(t, err)
		assert.Empty(t, profiles)
	})

	t.Run("Invalid file content", func(t *testing.T) {
		invalidPath := filepath.Join(tmpDir, "invalid")
		err := os.WriteFile(invalidPath, []byte("invalid ini content"), 0644)
		assert.NoError(t, err)

		profiles, err := parseProfilesFromFile(invalidPath, false)
		assert.Error(t, err)
		assert.Empty(t, profiles)
	})
}
