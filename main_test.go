package main

import (
	"strings"
	"testing"
)

func TestGetFileOK(t *testing.T) {
	f, e := getFile("./main.go")
	if e != nil {
		t.Errorf("got error %v", e)
	}

	if f.Name() != "./main.go" {
		t.Errorf("got %v want ./main.go", f.Name())
	}
}

func TestGetFileNoFile(t *testing.T) {
	f, e := getFile("./nofile")
	if e == nil {
		t.Errorf("expected error, go %v", f)
	}
}

func TestReadFileOK(t *testing.T) {
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

func TestMapResponsesOK(t *testing.T) {
	m, err := mapResponses("./example.json")
	if err != nil {
		t.Errorf("got %v", err)
	}

	if _, ok := m["/api/1"]; !ok {
		t.Errorf("got %v wanted true", ok)
	}
}

func TestMapResponsesNoFile(t *testing.T) {
	_, err := mapResponses("./nofile")
	if err == nil {
		t.Errorf("got %v expected an error", err)
	}
}

func TestMapResponsesBadJSON(t *testing.T) {
	m, err := mapResponses("./invalid.json")
	if err == nil {
		t.Errorf("got %v", err)
	}

	if _, ok := m["/api/1"]; ok {
		t.Errorf("got %v expected error", ok)
	}
}
