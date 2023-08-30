package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type mockResponse struct {
	Headers    map[string]string      `json:"headers"`
	StatusCode int                    `json:"statusCode"`
	Body       map[string]interface{} `json:"body"`
}

func main() {
	var (
		file string
		port int
	)

	flag.StringVar(&file, "f", "", "JSON file to be parsed")
	flag.IntVar(&port, "p", 8000, "Port to be used")
	flag.Parse()

	if file == "" {
		fmt.Fprint(os.Stderr, "No file specified. moco -f FILENAME [-p PORT]\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Starting server on port: %d\n", port)
	http.Handle("/", requestHandler(file))
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func makePatterns(routes map[string]mockResponse) []*regexp.Regexp {
	var out []*regexp.Regexp
	for k := range routes {
		r := regexp.MustCompile(k)
		out = append(out, r)
	}
	return out
}

func mapResponses(file string) (map[string]mockResponse, error) {
	f, err := getFile(file)
	if err != nil {
		return nil, err
	}

	c := readFile(f)

	m := make(map[string]mockResponse, 0)
	if err := json.Unmarshal([]byte(strings.Join(c, "")), &m); err != nil {
		return nil, err
	}

	return m, err
}

func getFile(path string) (*os.File, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	return f, nil
}

func readFile(f *os.File) []string {
	var (
		r *bufio.Reader
		s string
		e error
		t []string
	)
	r = bufio.NewReader(f)
	s, e = readln(r)

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

func requestHandler(file string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			routes map[string]mockResponse
			err    error
			regs   []*regexp.Regexp
			mr     mockResponse
			found  bool
		)

		routes, err = mapResponses(file)
		if err != nil {
			http.Error(w, "Error parsing JSON file", http.StatusInternalServerError)
			return
		}
		regs = makePatterns(routes)

		for i := 0; i < len(regs); i++ {
			match := regs[i].FindAllString(r.URL.Path, -1)
			if len(match) > 0 && !found {
				found = true
				mr = routes[regs[i].String()]
			}
		}

		if _, ok := routes[r.URL.Path]; ok && !found {
			found = true
			mr = routes[r.URL.Path]
		}

		if found {
			var sc int
			if mr.StatusCode > 0 {
				sc = mr.StatusCode
			} else {
				sc = http.StatusOK
			}

			if mr.Body == nil {
				w.WriteHeader(sc)
				return
			}

			enc := json.NewEncoder(w)
			if mr.Headers != nil {
				for k, v := range mr.Headers {
					w.Header().Set(k, v)
				}
			}

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

		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})
}
