// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"fmt"

	"github.com/openfaas/faas-cli/version"
	"github.com/spf13/cobra"
)

// GitCommit injected at build-time
var (
	shortVersion bool
)

func init() {
	versionCmd.Flags().BoolVar(&shortVersion, "short-version", false, "Just print Git SHA")

	faasCmd.AddCommand(versionCmd)
}

// versionCmd displays version information
var versionCmd = &cobra.Command{
	Use:   "version [--short-version]",
	Short: "Display the clients version information",
	Long: fmt.Sprintf(`The version command returns the current clients version information.

This currently consists of the GitSHA from which the client was built.
- https://github.com/openfaas/faas-cli/tree/%s`, version.GitCommit),
	Example: `  faas-cli version
  faas-cli version --short-version`,
	Run: runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	if shortVersion {
		fmt.Println(version.BuildVersion())
	} else {
		printFiglet()
		fmt.Printf("Commit: %s\n", version.GitCommit)
		fmt.Printf("Version: %s\n", version.BuildVersion())
	}
}
