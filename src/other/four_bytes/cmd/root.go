package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	fourbytes "github.com/theQRL/trueblocks-core/src/other/four_bytes/pkg/fourBytes"
)

var rootCmd = &cobra.Command{
	Use:   "four_bytes",
	Short: "Generates cross-product four-bytes database",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		outDir, err := cmd.Flags().GetString("outdir")
		if err != nil {
			return
		}
		if outDir == "" {
			return errors.New("outdir is required")
		}
		fmt.Printf("Generating database (with %d goroutines each step)\n", runtime.GOMAXPROCS(0))
		return fourbytes.GenerateFromFiles(outDir)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("outdir", "o", "", "Where to save chunks")
}
