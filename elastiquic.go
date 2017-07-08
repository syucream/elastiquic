package main

import (
  "fmt"
  "net/http"

  "github.com/devsisters/goquic"
)

func main() {
  client := &http.Client{
    Transport: goquic.NewRoundTripper(false),
  }
  resp, err := client.Get("https://www.google.co.jp/")

  fmt.Println(resp)
  fmt.Println(err)
}
