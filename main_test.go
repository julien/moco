package main

import "testing"

func TestGetFile(t *testing.T) {

	f, e := getFile("./main.go")
	if e != nil {
		t.Errorf("got %v", e)
	}

	if f.Name() != "./main.go" {
		t.Errorf("got %v want \"./main.go\"", f.Name())
	}
}
