package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"

  "github.com/spf13/cobra"
)

//LaunchServerCmd .
var LaunchServerCmd = &cobra.Command{
  Use:   "launchServer [host:port] [name] ",
  Short: "Launch a new cipher server for tests, should be accessible using host:port",
  Run: func(cmd *cobra.Command, args []string) {
    if err := cipherCli.launchServer(cmd, args); err != nil {
      fmt.Printf("Error: %v\n", err)
      os.Exit(1)
    }
  },
}

func init() {
  RootCmd.AddCommand(LaunchServerCmd)
}

func (m *cipherCLI) launchServer(cmd *cobra.Command, args []string) error {
  if len(args) < 2 {
    return fmt.Errorf("usage cipher3 launchServer [name] [host:port]")
  }
  list := strings.Split(args[0], ":")
  if len(list) != 2 {
    return fmt.Errorf("the first argument should have host:port format")
  }
  launchServer(args[1], list[0], list[1])
  return nil
}

func launchServer(name string, host string, port string) {
  if host == "" {
    host = "localhost"
  }
  exec := exec.Command("testServer", name, host, port)
  exec.Stdout = os.Stdout
  exec.Stderr = os.Stderr
  err := exec.Start()
  if err != nil {
    fmt.Printf("Server start error: %v\n", err)
  }
}
