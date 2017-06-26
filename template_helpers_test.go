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
		t.Fatalf("got: %s, wanted: %s", r, value)
	}
}
