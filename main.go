package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
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
	regs     []*regexp.Regexp
)

func init() {
	flag.StringVar(&fileFlag, "f", "", "JSON file to be parsed")
	flag.IntVar(&portFlag, "p", 8000, "Port to be used")
}

func main() {
	flag.Parse()

	if fileFlag == "" {
		color.Red("No file specified. moco -f FILENAME [-p PORT]")
		os.Exit(1)
	}

	var err error
	routes, err = mapResponses(fileFlag)
	if err != nil {
		color.Red("Unable to parse JSON file")
		os.Exit(1)
	}

	// regs := makePatterns(routes)

	color.Cyan("Starting server on port: %v\n", portFlag)
	http.Handle("/", requestHandler())
	http.ListenAndServe(":"+strconv.Itoa(portFlag), nil)
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

func requestHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// "reload" json file
		var err error
		routes, err = mapResponses(fileFlag)
		if err != nil {
			http.Error(w, "Error parsing JSON file", http.StatusInternalServerError)
			return
		}
		regs = makePatterns(routes)

		// check if path matches
		var mr mockResponse
		var found bool

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
