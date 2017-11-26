package main

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "time"

  "github.com/freignat91/cipher3/core"
  "github.com/spf13/cobra"
)

//CreateKeyCmd create the main key
var CreateKeyCmd = &cobra.Command{
  Use:   "createKey [keysPath] [dimension] [size]",
  Short: "Create cipher key",
  Run: func(cmd *cobra.Command, args []string) {
    if err := cipherCli.createCipherKey(cmd, args); err != nil {
      fmt.Printf("Error: %v\n", err)
      os.Exit(1)
    }
  },
}

func init() {
  RootCmd.AddCommand(CreateKeyCmd)
  CreateKeyCmd.Flags().BoolP("random", "r", false, "add manual random keybord inputs to ensure key is random")
}

func (m *cipherCLI) createCipherKey(cmd *cobra.Command, args []string) error {
  if len(args) < 3 {
    return fmt.Errorf("usage: cipher3 createKey [keysPath] [dimension] [size]")
  }
  dimension, err := strconv.Atoi(args[1])
  if err != nil {
    return fmt.Errorf("dimension is not a number")
  }
  keyBitSize, err := strconv.Atoi(args[2])
  if err != nil {
    return fmt.Errorf("size is not a number")
  }
  path := args[0]
  var randomList = make([]string, dimension, dimension)
  if cmd.Flag("random").Value.String() == "true" {
    getRandomList(&randomList)
  }
  t0 := time.Now()
  core, err := core.CreateKey(dimension, keyBitSize, randomList, m.verbose, m.debug)
  if err != nil {
    return err
  }
  if m.verbose {
    fmt.Printf("Compute time=%ds\n", time.Now().Sub(t0).Nanoseconds()/1000000000)
  }
  if err := core.SaveKey(path); err != nil {
    return err
  }
  return nil
}

func getRandomList(list *[]string) {
  for ii := range *list {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("Enter random keyboard input for dimension %d: ", ii+1)
    text, _ := reader.ReadString('\n')
    (*list)[ii] = text
  }
}
