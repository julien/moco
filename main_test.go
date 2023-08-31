package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

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
	f, e := getFile("./fixtures/fixtures.txt")
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
	m, err := mapResponses("./fixtures/example.json")
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
	m, err := mapResponses("./fixtures/invalid.json")
	if err == nil {
		t.Errorf("got %v expect error", err)
	}

	if _, ok := m["/api/1"]; ok {
		t.Errorf("got %v expected error", ok)
	}
}

func TestRequestHandlerOK(t *testing.T) {
	file := "./fixtures/example.json"
	_, err := mapResponses(file)
	if err != nil {
		t.Errorf("got %v", err)
	}

	handle := requestHandler(file)
	req, err := http.NewRequest("GET", "/api/1", nil)
	w := httptest.NewRecorder()

	handle.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 200", w.Code)
	}
}

func TestRequestHandlerNoBody(t *testing.T) {
	file := "./fixtures/nobody.json"
	handle := requestHandler(file)
	req, _ := http.NewRequest("GET", "/api/1", nil)
	w := httptest.NewRecorder()

	handle.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 403", w.Code)
	}
}

func TestRequestHandlerNoFile(t *testing.T) {
	handle := requestHandler("")
	req, _ := http.NewRequest("GET", "/api/1", nil)
	w := httptest.NewRecorder()

	handle.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %v want 200", w.Code)
	}
}

func TestRequestHandlerNotFound(t *testing.T) {
	file := "./fixtures/nobody.json"
	handle := requestHandler(file)
	req, _ := http.NewRequest("GET", "/api/2", nil)
	w := httptest.NewRecorder()

	handle.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("got %v want 404", w.Code)
	}
}
