package main

import (
	"os"
	"testing"
)

func TestConfigRequiresReleaseMetadata(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("RELEASE_SHA", "")
	if _, err := load(); err == nil {
		t.Fatal("expected validation error")
	}
	_ = os.Setenv("APP_ENV", "test")
	_ = os.Setenv("RELEASE_SHA", "abc")
	c, err := load()
	if err != nil || c.Environment != "test" {
		t.Fatalf("unexpected config: %+v %v", c, err)
	}
}
