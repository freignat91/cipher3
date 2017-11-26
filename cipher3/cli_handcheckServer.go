package main

import (
  "fmt"
  "os"
  "strings"

  "github.com/freignat91/cipher3/testServer/cipherserver"
  "github.com/spf13/cobra"
)

const localhost = "localhost"

//HandcheckServerCmd .
var HandcheckServerCmd = &cobra.Command{
  Use:   "handcheckServer [fromHost:fromPort] [toHost:toPort]",
  Short: "handcheck from one server to another",
  Run: func(cmd *cobra.Command, args []string) {
    if err := cipherCli.handcheckServer(cmd, args); err != nil {
      fmt.Printf("Error: %v\n", err)
      os.Exit(1)
    }
  },
}

func init() {
  RootCmd.AddCommand(HandcheckServerCmd)
}

func (m *cipherCLI) handcheckServer(cmd *cobra.Command, args []string) error {
  if len(args) < 2 {
    return fmt.Errorf("usage: cipher3 handcheckServer [fromHost:fromPort] [toHost:toPort]")
  }
  from := strings.Split(args[0], ":")
  if len(from) != 2 {
    return fmt.Errorf("the first argument should have host:port format")
  }
  if from[0] == "" {
    from[0] = localhost
  }
  to := strings.Split(args[1], ":")
  if len(to) != 2 {
    return fmt.Errorf("the second argument should have host:port format")
  }
  if to[0] == "" {
    to[0] = localhost
  }
  if err := cipherserver.Handcheck(from[0], from[1], to[0], to[1]); err != nil {
    return err
  }
  return nil
}
