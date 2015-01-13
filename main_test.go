package main

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
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
		t.Errorf("got %v expect error", err)
	}

	if _, ok := m["/api/1"]; ok {
		t.Errorf("got %v expected error", ok)
	}
}

func TestRequestHandlerKO(t *testing.T) {

	_, err := mapResponses("./example.json")
	if err != nil {
		t.Errorf("got %v", err)
	}

	handle := requestHandler()
	req, err := http.NewRequest("GET", "/api/1", nil)
	w := httptest.NewRecorder()

	handle.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %v want 500", w.Code)
	}
}

func TestRequestHandlerOK(t *testing.T) {

	_, err := mapResponses("./example.json")
	if err != nil {
		t.Errorf("got %v", err)
	}

	handle := requestHandler()
	fileFlag = "./example.json"
	req, err := http.NewRequest("GET", "/api/1", nil)
	w := httptest.NewRecorder()

	handle.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 200", w.Code)

	}
}

func TestRequestHandlerNoBody(t *testing.T) {

	handle := requestHandler()
	fileFlag = "./nobody.json"
	req, _ := http.NewRequest("GET", "/api/1", nil)
	w := httptest.NewRecorder()

	handle.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %v want 403", w.Code)

	}

	if w.Header()["Content-Type"][0] != "text/plain; charset=utf-8" {
		t.Errorf("got %v wanted text/plain; charset=utf-8", w.Header()["Content-Type"])
	}
}
