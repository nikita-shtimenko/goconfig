package env_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	goconfig "github.com/nikita-shtimenko/goconfig"
	"github.com/nikita-shtimenko/goconfig/loader/env"
)

type SampleConfig struct {
	AppName string `env:"APP_NAME"`
	Port    int    `env:"PORT"`
}

type testCase struct {
	name           string
	envContent     string
	envFiles       []string
	opts           []env.LoaderOption
	expectError    bool
	expectedConfig *SampleConfig
	errorContains  string
}

func TestLoader(t *testing.T) {
	tests := []testCase{
		{
			name: "Valid env file",
			envContent: `
				APP_NAME=testapp
				PORT=8080
			`,
			expectError: false,
			expectedConfig: &SampleConfig{
				AppName: "testapp",
				Port:    8080,
			},
		},
		{
			name:          "Missing env file error",
			envFiles:      []string{"nonexistent.env"},
			expectError:   true,
			errorContains: "error loading env file",
		},
		{
			name: "Skip missing env files",
			envContent: `
				APP_NAME=skippy
				PORT=3000
			`,
			opts: []env.LoaderOption{env.WithSkipMissingFiles()},
			envFiles: []string{
				"missing.env", // will be skipped
			},
			expectError: false,
			expectedConfig: &SampleConfig{
				AppName: "skippy",
				Port:    3000,
			},
		},
		{
			name:          "No env files",
			envFiles:      []string{},
			expectError:   true,
			errorContains: env.ErrEnvFilesNotSpecified.Error(),
		},
		{
			name:          "Invalid env values",
			envContent:    "PORT=notanumber",
			expectError:   true,
			errorContains: "parsing env vars",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer clearEnvironmentVariables("APP_NAME", "PORT")
			runTestCase(t, tc)
		})
	}
}

func runTestCase(t *testing.T, tc testCase) {
	envFiles := tc.envFiles
	if tc.envContent != "" {
		envFile := createTempEnvFile(t, tc.envContent)
		envFiles = append([]string{envFile}, envFiles...)
	}

	loader, err := env.NewLoader[SampleConfig](envFiles, tc.opts...)
	if tc.expectError {
		assertExpectedError(t, loader, err, tc)
		return
	}

	if err != nil {
		t.Fatalf("failed to create env loader: %v", err)
	}

	cfg, err := goconfig.NewConfig(loader)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	assertConfigValues(t, cfg, tc.expectedConfig)
}

func assertExpectedError[T any](t *testing.T, loader goconfig.ConfigLoader[T], err error, tc testCase) {
	if err == nil && tc.errorContains != "" {
		_, loadErr := goconfig.NewConfig(loader)
		if loadErr == nil {
			t.Fatal("expected config loading error, got nil")
		}
		if !strings.Contains(loadErr.Error(), tc.errorContains) {
			t.Errorf("expected error to contain '%s', got '%v'", tc.errorContains, loadErr)
		}
	} else if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
		t.Errorf("expected loader error to contain '%s', got '%v'", tc.errorContains, err)
	}
}

func assertConfigValues(t *testing.T, got *SampleConfig, want *SampleConfig) {
	t.Helper()
	if got.AppName != want.AppName {
		t.Errorf("AppName: expected %q, got %q", want.AppName, got.AppName)
	}
	if got.Port != want.Port {
		t.Errorf("Port: expected %d, got %d", want.Port, got.Port)
	}
}

func createTempEnvFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	file := filepath.Join(dir, ".env")

	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}

	return file
}

func clearEnvironmentVariables(keys ...string) {
	for _, k := range keys {
		_ = os.Unsetenv(k)
	}
}
