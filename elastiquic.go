package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/devsisters/goquic"
)

const (
	DEFINITIONS_FILE = "definitions.json"

	// error message templates
	STATUS_CODE_ERRMSG = "StatusCode: Expected is %d , but actual is %d.\n"
)

type Definitions struct {
	Scenarios []Scenario
}

type Scenario struct {
	Url     string
	Expects Expects
}

type Expects struct {
	StatusCode int
}

type TestResult struct {
	Successed    bool
	Url          string
	ErrorMessage string
}

type Stats struct {
	Successed int
	Failed    int
}

// Load json
func load() (Definitions, error) {
	file, err := ioutil.ReadFile(DEFINITIONS_FILE)
	if err != nil {
		return Definitions{}, err
	}

	var defs Definitions
	json.Unmarshal(file, &defs)

	return defs, nil
}

// Do QUIC request
func request(client *http.Client, scenario Scenario, ch chan TestResult) {
	resp, err := client.Get(scenario.Url)
	ch <- spec(scenario, resp, err)
}

// Check QUIC response
func spec(scenario Scenario, resp *http.Response, err error) TestResult {
	r := TestResult{true, scenario.Url, ""}
	expects := scenario.Expects

	if err != nil {
		r.Successed = false
		r.ErrorMessage = err.Error()
		return r
	}

	if expects.StatusCode != 0 && expects.StatusCode != resp.StatusCode {
		r.Successed = false
		r.ErrorMessage = fmt.Sprintf(STATUS_CODE_ERRMSG, expects.StatusCode, resp.StatusCode)
		return r
	}

	return r
}

func main() {
	defs, err := load()
	if err != nil {
		fmt.Println("elastiquic can't load a JSON file.")
		os.Exit(1)
	}

	client := &http.Client{
		Transport: goquic.NewRoundTripper(false),
	}

	// TODO Controll concurrency
	ch := make(chan TestResult)
	for _, scenario := range defs.Scenarios {
		go request(client, scenario, ch)
	}

	stats := Stats{0, 0}
	for range defs.Scenarios {
		result := <-ch

		if result.Successed {
			stats.Successed += 1
			fmt.Print(".")
		} else {
			stats.Failed += 1
			fmt.Printf("X%s is failed because %s\n", result.Url, result.ErrorMessage)
		}
	}

	total := stats.Successed + stats.Failed
	fmt.Printf("\n\nTotal requests: %d, successed: %d, failed: %d\n", total, stats.Successed, stats.Failed)
}
