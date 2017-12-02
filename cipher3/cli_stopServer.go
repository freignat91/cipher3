package main

import (
  "fmt"
  "os"
  "strings"

  "github.com/freignat91/cipher3/testServer/cipherserver"
  "github.com/spf13/cobra"
)

//StopServerCmd .
var StopServerCmd = &cobra.Command{
  Use:   "stopServer [host:port]",
  Short: "stop a cipher server",
  Run: func(cmd *cobra.Command, args []string) {
    if err := cipherCli.stopServer(cmd, args); err != nil {
      fmt.Printf("Error: %v\n", err)
      os.Exit(1)
    }
  },
}

func init() {
  RootCmd.AddCommand(StopServerCmd)
}

func (m *cipherCLI) stopServer(cmd *cobra.Command, args []string) error {
  if len(args) < 1 {
    return fmt.Errorf("usage: cipher3 stopServer [host:port]")
  }
  list := strings.Split(args[0], ":")
  if len(list) != 2 {
    return fmt.Errorf("the first argument should have host:port format")
  }
  if list[0] == "" {
    list[0] = "localhost"
  }
  cipherserver.StopCipherServer(list[0], list[1])
  return nil
}
