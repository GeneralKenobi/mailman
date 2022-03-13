package config

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	// This instance is used for the tests so that modifying the default configuration in default.go doesn't break the tests.
	mockDefaultConfig := Config{
		HttpServer: HttpServer{
			Port:                   8080,
			ShutdownTimeoutSeconds: 30,
		},
	}
	tests := map[string]struct {
		mockFileReadFunc     func(filepath string) ([]byte, error)
		inputConfigFilePaths []string
		expected             Config
		expectError          bool
	}{
		"Should merge the default config with the configs loaded from files": {
			mockFileReadFunc: func(filepath string) ([]byte, error) {
				switch filepath {
				case "/first/config.json":
					return []byte(`{"httpServer": {"port":8443}, "postgres": {"port": 5432}}`), nil
				case "/second/config.json":
					return []byte(`{"httpServer": {"port":9099}}`), nil
				default:
					return nil, fmt.Errorf("file %q doesn't exist", filepath)
				}
			},
			inputConfigFilePaths: []string{"/first/config.json", "/second/config.json"},
			expected: Config{
				HttpServer: HttpServer{
					Port:                   9099, // From second custom config file
					ShutdownTimeoutSeconds: 30,   // From defaults
				},
				Postgres: Postgres{
					Port: 5432, // From first custom config file
				},
			},
			expectError: false,
		},
		"Should return the default config if config file path is empty": {
			mockFileReadFunc: func(filepath string) ([]byte, error) {
				t.Errorf("Shouldn't have been called")
				return nil, fmt.Errorf("shouldn't have been called")
			},
			inputConfigFilePaths: nil,
			expected:             mockDefaultConfig,
			expectError:          false,
		},
	}

	originalDefaultConfiguration := defaultConfig
	originalFileReadHook := fileReadHook
	defer func() {
		defaultConfig = originalDefaultConfiguration
		fileReadHook = originalFileReadHook
	}()

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			defaultConfig = mockDefaultConfig
			fileReadHook = test.mockFileReadFunc

			result, err := Load(test.inputConfigFilePaths)

			// TODO: Extract error comparison to helper func
			// TODO: Extract pretty print to helper func
			if test.expectError && err == nil {
				t.Errorf("Expected an error but got none")
			}
			if !test.expectError && err != nil {
				t.Errorf("Expected no error but got %v", err)
			}
			if result != test.expected {
				t.Errorf("Expected %+v\nGot %+v", test.expected, result)
			}
		})
	}
}
