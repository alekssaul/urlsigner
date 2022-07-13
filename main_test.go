package main

import (
	"testing"
)

func TestStringToUnix(t *testing.T) {
	_, err := stringToUnix("foo")
	if err == nil {
		t.Errorf("Function stringToUnix did not error when called with \"foo\" ")
	}

	// 1657680270 - Wednesday, July 13, 2022 2:44:30 AM
	a, err := stringToUnix("1657680270")
	if err != nil {
		t.Errorf("Function stringToUnix errored when called with \"1657680270\" ")
	}

	if a.Unix() != 1657680270 {
		t.Errorf("Function stringToUnix errored when called with \"1657680270\" ")
	}
}

