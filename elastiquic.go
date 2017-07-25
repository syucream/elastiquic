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
  Url string
  Expects Expects
}

type Expects struct {
  StatusCode int
}

type TestResult struct {
  Successed bool
  Url string
  ErrorMessage string
}

// Load json
func load() Definitions {
  file, err := ioutil.ReadFile(DEFINITIONS_FILE)
  if err != nil {
    fmt.Println(err)
    // FIXME
    os.Exit(1)
  }

  var defs Definitions
  json.Unmarshal(file, &defs)

  return defs
}

// Do QUIC request
func request(client *http.Client, scenario Scenario, ch chan TestResult) {
  resp, err := client.Get(scenario.Url)

  // For debug
  // fmt.Println(resp)
  // fmt.Println(err)

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

  if expects.StatusCode != 0 {
    if expects.StatusCode != resp.StatusCode {
      r.Successed = false
      r.ErrorMessage = fmt.Sprintf(STATUS_CODE_ERRMSG, expects.StatusCode, resp.StatusCode)
      return r
    }
  }

  return r
}

func main() {
  defs := load()

  client := &http.Client {
    Transport: goquic.NewRoundTripper(false),
  }

  ch := make(chan TestResult)
  for _, scenario := range defs.Scenarios {
    go request(client, scenario, ch)
  }

  for range defs.Scenarios {
    result := <- ch
    if !result.Successed {
      fmt.Println(result.ErrorMessage)
    }
  }
}
