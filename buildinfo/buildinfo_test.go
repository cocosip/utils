package buildinfo

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	// Save original ldDevMode value and restore it after test
	originalLdDevMode := ldDevMode
	defer func() {
		ldDevMode = originalLdDevMode
		// Re-initialize info struct after changing ldDevMode for subsequent tests
		// This is a hack for testing ldflags-set variables, not typical usage.
		// In real applications, ldflags are set once at compile time.
		// We'll simulate it by re-assigning and then calling a helper to re-populate 'info'.
		// For simplicity, we'll just re-assign the global 'info' directly here.
		info = BuildInfo{
			Version:     ldVersion,
			BuildCommit: ldBuildCommit,
			BuildTime:   ldBuildTime,
			DevMode:     parseBool(ldDevMode),
		}
	}()

	bi := Get()

	// Test default values (unless ldflags are set during compilation)
	if bi.Version == "" {
		t.Errorf("Version should not be empty")
	}
	if bi.BuildCommit == "" {
		t.Errorf("BuildCommit should not be empty")
	}
	if bi.BuildTime == "" {
		t.Errorf("BuildTime should not be empty")
	}

	// Check if default values are present
	if bi.Version != "DEVELOPMENT Version" {
		t.Logf("Version: %s (expected DEVELOPMENT Version unless ldflags set)", bi.Version)
	}
	if bi.BuildCommit != "DEVELOPMENT Build" {
		t.Logf("BuildCommit: %s (expected DEVELOPMENT Build unless ldflags set)", bi.BuildCommit)
	}
	if bi.BuildTime != "UNKNOWN" {
		t.Logf("BuildTime: %s (expected UNKNOWN unless ldflags set)", bi.BuildTime)
	}

	// Test with ldDevMode = "true"
	ldDevMode = "true"
	info = BuildInfo{
		Version:     ldVersion,
		BuildCommit: ldBuildCommit,
		BuildTime:   ldBuildTime,
		DevMode:     parseBool(ldDevMode),
	}
	if !Get().DevMode {
		t.Errorf("DevMode should be true for 'true'")
	}

	// Test with ldDevMode = "false"
	ldDevMode = "false"
	info = BuildInfo{
		Version:     ldVersion,
		BuildCommit: ldBuildCommit,
		BuildTime:   ldBuildTime,
		DevMode:     parseBool(ldDevMode),
	}
	if Get().DevMode {
		t.Errorf("DevMode should be false for 'false'")
	}

	// Test with invalid ldDevMode string
	ldDevMode = "invalid"
	info = BuildInfo{
		Version:     ldVersion,
		BuildCommit: ldBuildCommit,
		BuildTime:   ldBuildTime,
		DevMode:     parseBool(ldDevMode),
	}
	if Get().DevMode {
		t.Errorf("DevMode should be false for 'invalid'")
	}
}

func TestBuildInfo_String(t *testing.T) {
	// Temporarily set ldflags variables for testing String() method
	originalLdVersion := ldVersion
	originalLdBuildCommit := ldBuildCommit
	originalLdBuildTime := ldBuildTime
	originalLdDevMode := ldDevMode
	defer func() {
		ldVersion = originalLdVersion
		ldBuildCommit = originalLdBuildCommit
		ldBuildTime = originalLdBuildTime
		ldDevMode = originalLdDevMode
		// Reset global 'info' struct for other tests
		info = BuildInfo{
			Version:     ldVersion,
			BuildCommit: ldBuildCommit,
			BuildTime:   ldBuildTime,
			DevMode:     parseBool(ldDevMode),
		}
	}()

	ldVersion = "v1.2.3"
	ldBuildCommit = "abcde123"
	ldBuildTime = "2023-01-01T12:00:00Z"
	ldDevMode = "false"
	info = BuildInfo{
		Version:     ldVersion,
		BuildCommit: ldBuildCommit,
		BuildTime:   ldBuildTime,
		DevMode:     parseBool(ldDevMode),
	}

	bi := Get()
	expectedString := "Version: v1.2.3\nBuild: abcde123\nBuilt At: 2023-01-01T12:00:00Z\nDev Mode: false"
	if bi.String() != expectedString {
		t.Errorf("String() mismatch\nExpected:\n%s\nGot:\n%s", expectedString, bi.String())
	}
}

func TestBuildInfo_JSON(t *testing.T) {
	// Temporarily set ldflags variables for testing JSON() method
	originalLdVersion := ldVersion
	originalLdBuildCommit := ldBuildCommit
	originalLdBuildTime := ldBuildTime
	originalLdDevMode := ldDevMode
	defer func() {
		ldVersion = originalLdVersion
		ldBuildCommit = originalLdBuildCommit
		ldBuildTime = originalLdBuildTime
		ldDevMode = originalLdDevMode
		// Reset global 'info' struct for other tests
		info = BuildInfo{
			Version:     ldVersion,
			BuildCommit: ldBuildCommit,
			BuildTime:   ldBuildTime,
			DevMode:     parseBool(ldDevMode),
		}
	}()

	ldVersion = "v1.2.3"
	ldBuildCommit = "abcde123"
	ldBuildTime = "2023-01-01T12:00:00Z"
	ldDevMode = "true"
	info = BuildInfo{
		Version:     ldVersion,
		BuildCommit: ldBuildCommit,
		BuildTime:   ldBuildTime,
		DevMode:     parseBool(ldDevMode),
	}

	bi := Get()
	jsonString, err := bi.JSON()
	if err != nil {
		t.Fatalf("JSON() returned an error: %v", err)
	}

	expectedJSON := `{
  "version": "v1.2.3",
  "buildCommit": "abcde123",
  "buildTime": "2023-01-01T12:00:00Z",
  "devMode": true
}`

	// Compare JSON strings, ignoring whitespace differences for robustness
	if strings.TrimSpace(jsonString) != strings.TrimSpace(expectedJSON) {
		t.Errorf("JSON() mismatch\nExpected:\n%s\nGot:\n%s", expectedJSON, jsonString)
	}

	// Test unmarshaling to ensure JSON is valid
	var unmarshaledBi BuildInfo
	err = json.Unmarshal([]byte(jsonString), &unmarshaledBi)
	if err != nil {
		t.Fatalf("Failed to unmarshal generated JSON: %v", err)
	}

	if unmarshaledBi.Version != bi.Version ||
		unmarshaledBi.BuildCommit != bi.BuildCommit ||
		unmarshaledBi.BuildTime != bi.BuildTime ||
		unmarshaledBi.DevMode != bi.DevMode {
		t.Errorf("Unmarshaled JSON does not match original BuildInfo")
	}
}
