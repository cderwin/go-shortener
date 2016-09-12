package main

import (
	"testing"
	"time"
)

// Mock Clock

var MockNow = time.Date(2016, time.January, 168, 0, 0, 0, 0, time.UTC)

type MockClock struct {
	current time.Time
}

func CreateMockClock() MockClock {
	return MockClock{current: MockNow}
}

func (c MockClock) UTCNow() time.Time {
	return c.current
}


// Actual tests

func TestBase62(t *testing.T) {
	expectedMap := map[uint32]string{14: "o", 62: "ba", 3843: "99", 1569087: "gKl1", 384: "gm"}
	for value, expected := range expectedMap {
		actual := base62Encode(value)
		if actual != expected {
			t.Errorf("Input int: %v\nExpected string: %s\nActual string: %s", value, expected, actual)
		}
	}
}

func TestHashURL(t *testing.T) {
	expectedMap := map[string]string{"hello": "9x58c", "goodbye": "pyRkS", "welcome to fallujah": "erw1HB"}
	for value, expected := range expectedMap {
		actual := hashUrl(value)
		if actual != expected {
			t.Errorf("Input string: %s\nExpected string: %s\nActual string: %s", value, expected, actual)
		}
	}
}
