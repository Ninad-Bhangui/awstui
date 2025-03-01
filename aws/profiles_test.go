package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProfiles(t *testing.T) {
	// Set up test AWS credentials file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".aws", "config")
	os.MkdirAll(filepath.Dir(configPath), 0755)

	// Corrected config file content with "profile" prefix
	configContent := `[profile default]
region = us-east-1

[profile dev]
region = us-west-2
`
	fmt.Println(configPath)
	os.WriteFile(configPath, []byte(configContent), 0644)

	// Test ParseProfilesFromConfig()
	profiles, err := ParseProfilesFromConfig(configPath)
	fmt.Println(err)
	assert.NoError(t, err)
	assert.Contains(t, profiles, "default")
	assert.Contains(t, profiles, "dev")
}
