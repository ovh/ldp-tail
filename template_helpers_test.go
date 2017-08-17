package main

import (
	"os"
	"testing"
)

func TestDate(t *testing.T) {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	r := date(1498508631)
	value := "2017-06-26 20:23:51"
	if r != value {
		t.Fatalf("Date with default format, got: %s, wanted: %s", r, value)
	}

	r = date(1498508631, "02/01/2006:15:04:05")
	value = "26/06/2017:20:23:51"
	if r != value {
		t.Fatalf("Date with custom format, got: %s, wanted: %s", r, value)
	}
}

func TestLevel(t *testing.T) {
	r, err := level(3)
	value := "err"
	if err != nil {
		t.Fatalf("Level with int failure: %s", err.Error())
	}
	if r != value {
		t.Fatalf("Level with int, got: %s, wanted: %s", r, value)
	}

	r, err = level("4")
	value = "warn"
	if err != nil {
		t.Fatalf("Level with string failure: %s", err.Error())
	}
	if r != value {
		t.Fatalf("Level with string, got: %s, wanted: %s", r, value)
	}

	_, err = level("foo")
	if err == nil {
		t.Fatalf("Level with bad string should have returned an error")
	}
}
