package main

import (
	"fmt"
	"io"
	"os"

	"github.com/divingbeetle/jt/core"
	"github.com/spf13/cobra"
)

var collation string

var rootCmd = &cobra.Command{
	Use:   "jt [file]",
	Short: "Convert JSON to MySQL JSON_TABLE format",
	Long:  `A CLI tool that converts JSON data to MySQL JSON_TABLE format.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte
		var err error

		if len(args) > 0 {
			input, err = os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}
		} else {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
		}

		opts := core.Options{
			StringCollation: collation,
		}

		result, err := core.Convert(input, opts)
		if err != nil {
			return err
		}

		fmt.Println(result)
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&collation, "collation", "c", "", "collation for string columns")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
