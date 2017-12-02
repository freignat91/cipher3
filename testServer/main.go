package main

import (
  "fmt"
  "log"
  "os"

  "github.com/freignat91/cipher3/testServer/cipherserver"
)

// build vars
var (
  Version string
  Build   string
)

func main() {
  if len(os.Args) < 4 {
    fmt.Printf("Usage: testServer name host port\n")
    os.Exit(1)
  }
  server, err := cipherserver.NewCipherServer(os.Args[1], os.Args[2], os.Args[3])
  if err == nil {
    err := server.Start(Version)
    if err != nil {
      log.Printf("Exit on error: %v\n", err)
      os.Exit(1)
    }
  }
}
