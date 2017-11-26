package main

import (
  "fmt"
  "os"

  "github.com/freignat91/cipher3/testServer/cipherserver"
  "github.com/spf13/cobra"
)

//SendDataCmd .
var SendDataCmd = &cobra.Command{
  Use:   "sendData [host:port] [targetName] [data]",
  Short: "send encrypted [data] from [host:port] to the sever called [targetName]",
  Run: func(cmd *cobra.Command, args []string) {
    if err := cipherCli.sendData(cmd, args); err != nil {
      fmt.Printf("Error: %v\n", err)
      os.Exit(1)
    }
  },
}

func init() {
  RootCmd.AddCommand(SendDataCmd)
}

func (m *cipherCLI) sendData(cmd *cobra.Command, args []string) error {
  if len(args) < 3 {
    return fmt.Errorf("usage: cisendData [host:port] [targetName] [data]")
  }
  if err := cipherserver.SendData(args[0], args[1], args[2]); err != nil {
    return err
  }
  return nil
}
