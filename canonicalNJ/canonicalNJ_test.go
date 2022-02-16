package main

import "testing"

func TestFuture(t *testing.T) {
	goat := "monkey"
	swamp := "peanut"
	if swamp == goat {
		t.Errorf("hehehe")
	}
}
