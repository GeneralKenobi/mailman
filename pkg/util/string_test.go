package util

import (
	"regexp"
	"testing"
)

func TestRandomAlphanumericString(t *testing.T) {
	const targetLength = 16
	alphanumericMatcher := regexp.MustCompile("^[[:alnum:]]*$").MatchString

	result := RandomAlphanumericString(targetLength)

	if len(result) != targetLength {
		t.Errorf("Wanted to get %d characters but got %d in %v", targetLength, len(result), result)
	}
	if !alphanumericMatcher(result) {
		t.Errorf("Expected only alphanumeric characters in %v", result)
	}
}

func TestRandomAlphanumericStringShouldReturnEmptyStringForLengthEqualTo0(t *testing.T) {
	result := RandomAlphanumericString(0)

	if result != "" {
		t.Errorf("Expected an empty string for target length 0 but got %v", result)
	}
}

func TestRandomAlphanumericStringShouldPanicForNegativeLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for negative input")
		}
	}()

	RandomAlphanumericString(-1)
}
