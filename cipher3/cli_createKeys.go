package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/freignat91/cipher3/core"
	"github.com/spf13/cobra"
)

//CreateKeysCmd create the main key
var CreateKeysCmd = &cobra.Command{
	Use:   "createKeys [keysPath/name]",
	Short: "Create cipher keys",
	Long:  `Create cipher keys`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cipherCli.createCipherKey(cmd, args); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(CreateKeysCmd)
	CreateKeysCmd.Flags().StringP("dimension", "d", "3", "number of dimension: default 3")
	CreateKeysCmd.Flags().StringP("size", "s", "16384", "dimension size: should be a multiple of 64, at least 16384, default 16384")
}

func (m *cipherCLI) createCipherKey(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("need key file path as argument. usage: cipher createKeys --di--size [size]")
	}
	dimension, err := strconv.Atoi(cmd.Flag("dimension").Value.String())
	if err != nil {
		return fmt.Errorf("option --dimension is not a number")
	}
	keyBitSize, err := strconv.Atoi(cmd.Flag("size").Value.String())
	if err != nil {
		return fmt.Errorf("option --size is not a number")
	}
	path := args[0]
	t0 := time.Now()
	core, err := core.CreateKey(dimension, keyBitSize, m.verbose, m.debug)
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
