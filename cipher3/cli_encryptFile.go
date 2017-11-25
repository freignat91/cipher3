package main

import (
	"fmt"
	"time"

	"github.com/freignat91/cipher3/core"
	"github.com/spf13/cobra"
)

//EncryptFileCmd .
var EncryptFileCmd = &cobra.Command{
	Use:   "encryptFile [sourcefilePath] [targetFilePath] [keyFilePath]",
	Short: "encrypt file",
	Long:  `encrypt file`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cipherCli.encryptFile(cmd, args); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(EncryptFileCmd)
}

func (m *cipherCLI) encryptFile(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage cipher encryptFile [sourcefilePath] [targetFilePath] [keyFilePath]")
	}
	t0 := time.Now()
	if err := core.CipherFile(args[0], args[1], args[2]); err != nil {
		return err
	}
	fmt.Printf("done time=%ds\n", time.Now().Sub(t0).Nanoseconds()/1000000000)
	return nil
}
