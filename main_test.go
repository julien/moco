package main

import (
	"strings"
	"testing"
)

func TestGetFile(t *testing.T) {

	f, e := getFile("./main.go")
	if e != nil {
		t.Errorf("got error %v", e)
	}

	if f.Name() != "./main.go" {
		t.Errorf("got %v want ./main.go", f.Name())
	}

	// Test for a non existing file
	f, e = getFile("./nofile")
	if e == nil {
		t.Errorf("expected error, go %v", f)
	}
}

func TestReadFile(t *testing.T) {
	f, e := getFile("./fixtures.txt")
	if e != nil {
		t.Errorf("got error %v", e)
	}

	s := readFile(f)
	r := strings.Join(s, "")

	if r != "This is a file" {
		t.Errorf("got %v want This is a test", r)
	}
}
