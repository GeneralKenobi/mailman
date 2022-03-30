package mdctx

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestLogging(t *testing.T) {
	originalOutput := logger.Writer()
	originalCurrentTimeHook := currentTime
	defer func() {
		logger.SetOutput(originalOutput)
		currentTime = originalCurrentTimeHook
	}()

	var logOutput bytes.Buffer
	logger.SetOutput(&logOutput)
	currentTime = func() time.Time {
		return time.Date(2022, 3, 13, 20, 59, 48, 162428129, time.UTC)
	}

	ctx := context.Background()
	ctx = WithCorrelationId(ctx, "abc123")
	ctx = WithRequestUri(ctx, "/health")

	Infof(ctx, "Message with format: %s, %d", "value", 123)
	Warnf(ctx, "Warn message")
	Errorf(nil, "Error message without context")

	// Line numbers have to match the actual line numbers for the function calls above
	const expected = `
2022-03-13T20:59:48.162Z INFO                mdctx_test.go:29  ]   Message with format: value, 123   ][correlation-id=abc123][uri=/health]
2022-03-13T20:59:48.162Z WARN                mdctx_test.go:30  ]   Warn message   ][correlation-id=abc123][uri=/health]
2022-03-13T20:59:48.162Z ERROR               mdctx_test.go:31  ]   Error message without context   ]
`
	expectedTrimmed := strings.TrimSpace(expected)

	actual := logOutput.String()
	actualTrimmed := strings.TrimSpace(actual)

	if actualTrimmed != expectedTrimmed {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedTrimmed, actualTrimmed)
	}
}

func TestCorrelationId(t *testing.T) {
	tests := map[string]struct {
		input    context.Context
		expected string
	}{
		"Should return correlation ID if both correlation ID and operation ID are set": {
			input:    withValue(withValue(context.Background(), operationIdKey, "op-1"), correlationIdKey, "cor-1"),
			expected: "cor-1",
		},
		"Should return operation ID if operation ID is set but correlation ID is not set": {
			input:    withValue(context.Background(), operationIdKey, "op-1"),
			expected: "op-1",
		},
		"Should return empty string if neither correlation ID nor operation ID are set": {
			input:    context.Background(),
			expected: "",
		},
	}

	for title, test := range tests {
		t.Run(title, func(t *testing.T) {
			result := CorrelationId(test.input)
			if result != test.expected {
				t.Errorf("Expected %v but got %v", test.expected, result)
			}
		})
	}
}
