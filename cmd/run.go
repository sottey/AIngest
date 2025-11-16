/*
Copyright © 2025 sottey

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sottey/aingest/internal/bundler"
	"github.com/spf13/cobra"
)

var (
	outputFile string
	tokenMode  string
	recursive  bool
	format     string
)

var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Ingest and bundle files into structured output",
	Long: `Scans one or more input files, detects their format, 
normalizes them into a unified AI-ready bundle, and outputs structured JSON, XML, or SHON.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := args[0]
		fmt.Printf("🧠 AIngest started on: %s\n", root)

		b := bundler.NewBundler(root, recursive)
		bundle, err := b.BuildBundle()
		if err != nil {
			fmt.Println("Error building bundle:", err)
			os.Exit(1)
		}

		data, err := json.MarshalIndent(bundle, "", "  ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			os.Exit(1)
		}

		if err := os.WriteFile(outputFile, data, 0644); err != nil {
			fmt.Println("Error writing output:", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Bundle created successfully: %s\n", outputFile)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&outputFile, "out", "o", "bundle.json", "Output file path")
	runCmd.Flags().StringVarP(&format, "format", "f", "json", "Output format (json, xml, shon)")
	runCmd.Flags().StringVar(&tokenMode, "token-mode", "fast", "Token mode: fast or precise")
	runCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively scan subdirectories")
}
