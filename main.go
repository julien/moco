package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/fsnotify.v1"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type mockResponse struct {
	Headers    map[string]string      `json:"headers"`
	StatusCode int                    `json:"statusCode"`
	Body       map[string]interface{} `json:"body"`
}

var (
	fileFlag string
	portFlag int
	routes   map[string]mockResponse
	watcher  *fsnotify.Watcher
)

func init() {
	flag.StringVar(&fileFlag, "f", "", "JSON file to be parsed")
	flag.IntVar(&portFlag, "p", 8000, "Port to be used")
}

func main() {
	flag.Parse()

	if fileFlag == "" {
		fmt.Println("No file specified. [moco -f FILENAME]")
		os.Exit(1)
	}

	var err error
	routes, err = mapResponses(fileFlag)
	if err != nil {
		fmt.Println("Unable to parse JSON file")
		os.Exit(1)
	}

	go createWatcher()

	// Start server
	fmt.Printf("Starting server on port: %v\n", portFlag)
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":"+strconv.Itoa(portFlag), nil)
}

func mapResponses(file string) (map[string]mockResponse, error) {
	f, err := getFile(file)
	if err != nil {
		fmt.Println("Error reading file", err)
		return nil, err
	}

	c := readFile(f)

	m := make(map[string]mockResponse, 0)
	if err := json.Unmarshal([]byte(strings.Join(c, "")), &m); err != nil {
		fmt.Printf("Error parsing JSON: %v", err)
		return nil, err
	}

	return m, err
}

func createWatcher() (*fsnotify.Watcher, error) {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Watcher create error %s\n", err)
		return nil, err
	}
	defer watcher.Close()

	if err = watcher.Add(fileFlag); err != nil {
		fmt.Printf("Watcher add error %s\n", err)
		return nil, err
	}

	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf("JSON file changed: %s\n", ev.Name)

				routes, _ = mapResponses(fileFlag)

			}
		case err := <-watcher.Errors:
			fmt.Printf("Watch error %s:\n", err)
		}
	}

	return watcher, nil
}

func getFile(path string) (*os.File, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	return f, nil
}

func readFile(f *os.File) []string {

	r := bufio.NewReader(f)
	s, e := readln(r)
	var t []string

	for e == nil {
		s = strings.Trim(s, " \n\t\r")
		if s != "" && len(s) > 0 {
			t = append(t, s)
		}
		s, e = readln(r)
	}
	return t
}

func readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

	if mr, ok := routes[r.URL.Path]; ok {

		if mr.Body == nil {
			http.Error(w, "No response body defined for this request", http.StatusBadRequest)
			return
		}

		enc := json.NewEncoder(w)
		if mr.Headers != nil {
			for k, v := range mr.Headers {
				w.Header().Set(k, v)
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
		}

		// Cache headers
		age := 30 * 24 * 60 * 60 * 1000
		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(age))
		t := time.Now().Add(time.Duration(time.Hour*24) * 30)
		w.Header().Set("Expires", t.Format(time.RFC1123Z))

		if mr.StatusCode != 0 {
			w.WriteHeader(mr.StatusCode)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		enc.Encode(mr.Body)
	}

}
