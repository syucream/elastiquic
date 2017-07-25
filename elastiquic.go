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
func request(client *http.Client, scenario Scenario, ch chan int) {
  resp, err := client.Get(scenario.Url)

  // For debug
  // fmt.Println(resp)
  // fmt.Println(err)

  if err != nil {
    ch <- 0
  } else {
    spec(scenario, resp)
    ch <- 1
  }
}

func spec(scenario Scenario, resp *http.Response) TestResult {
  r := TestResult{true, scenario.Url, ""}
  expects := scenario.Expects

  if expects.StatusCode != 0 {
    successed := expects.StatusCode == resp.StatusCode
    if !successed {
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

  ch := make(chan int)
  go request(client, defs.Scenarios[0], ch)
  <- ch
}
