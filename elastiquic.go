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
)

type Definitions struct {
  Scenarios []Scenario
}

type Scenario struct {
  Url string
  Expects Expects
}

type Expects struct {
  Statuscode int
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
  fmt.Println(resp)
  fmt.Println(err)

  ch <- 1
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
