package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/devsisters/goquic"
)

const (
	DEFINITIONS_FILE = "definitions.json"

	// error message templates
	STATUS_CODE_ERRMSG = "StatusCode: Expected is %d , but actual is %d.\n"
)

type Definitions struct {
	Maxprocs  int
	Scenarios []Scenario
}

type Scenario struct {
	Url     string
	Expects map[string]interface{}
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
func request(client *http.Client, scenario Scenario, result *TestResult, done *sync.WaitGroup) {
	resp, err := client.Get(scenario.Url)
	spec(scenario, resp, err, result)
	done.Done()
}

// Check QUIC response
func spec(scenario Scenario, resp *http.Response, err error, result *TestResult) {
	expects := scenario.Expects

	result.Successed = true
	result.Url = scenario.Url
	result.ErrorMessage = ""

	if err != nil {
		result.Successed = false
		result.ErrorMessage = err.Error()
		return
	}

	if val, ok := expects["statuscode"].(float64); ok {
		expectedStatusCode := int(val)
		if resp.StatusCode != expectedStatusCode {
			result.Successed = false
			result.ErrorMessage = fmt.Sprintf(STATUS_CODE_ERRMSG, expectedStatusCode, resp.StatusCode)
			return
		}
	}
}

func printResults(results []TestResult) {
	stats := Stats{0, 0}

	for _, result := range results {
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

func main() {
	defs, err := load()
	if err != nil {
		fmt.Println("elastiquic can't load a JSON file.")
		os.Exit(1)
	}

	// Set concurrency
	maxprocs := defs.Maxprocs
	if os.Getenv("GOMAXPROCS") == "" && maxprocs > 0 {
		runtime.GOMAXPROCS(maxprocs)
	}

	client := &http.Client{
		Transport: goquic.NewRoundTripper(false),
	}

	num := len(defs.Scenarios)
	var done sync.WaitGroup
	done.Add(num)

	// Do requests
	results := make([]TestResult, num)
	for i, scenario := range defs.Scenarios {
		go request(client, scenario, &results[i], &done)
	}
	done.Wait()

	printResults(results)
}
