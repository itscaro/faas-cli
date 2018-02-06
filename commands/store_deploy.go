// Copyright (c) OpenFaaS Project 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// Setup flags that are used only by deploy command (variables defined above)
	storeDeployCmd.Flags().StringArrayVarP(&storeDeployFlags.envvarOpts, "env", "e", []string{}, "Adds one or more environment variables to the defined ones by store (ENVVAR=VALUE)")
	storeDeployCmd.Flags().StringArrayVarP(&storeDeployFlags.labelOpts, "label", "l", []string{}, "Set one or more label (LABEL=VALUE)")
	storeDeployCmd.Flags().BoolVar(&storeDeployFlags.replace, "replace", false, "Replace any existing function")
	storeDeployCmd.Flags().BoolVar(&storeDeployFlags.update, "update", true, "Update existing functions")
	storeDeployCmd.Flags().StringArrayVar(&storeDeployFlags.constraints, "constraint", []string{}, "Apply a constraint to the function")
	storeDeployCmd.Flags().StringArrayVar(&storeDeployFlags.secrets, "secret", []string{}, "Give the function access to a secure secret")

	// Set bash-completion.
	_ = storeDeployCmd.Flags().SetAnnotation("handler", cobra.BashCompSubdirsInDir, []string{})

	storeCmd.AddCommand(storeDeployCmd)
}

var storeDeployCmd = &cobra.Command{
	Use: `deploy (FUNCTION_NAME|FUNCTION_TITLE)
                        [--gateway GATEWAY_URL]
                        [--network NETWORK_NAME]
                        [--env ENVVAR=VALUE ...]
                        [--label LABEL=VALUE ...]
                        [--replace=false]
                        [--update=true]
                        [--constraint PLACEMENT_CONSTRAINT ...]
                        [--secret "SECRET_NAME"]
                        [--url STORE_URL]`,

	Short: "Deploy OpenFaaS functions from a store",
	Long:  `Same as faas-cli deploy except that function is pre-loaded with arguments from the store`,
	Example: `  faas-cli store deploy figlet
  faas-cli store deploy figlet \
    --gateway=http://localhost:8080 \
    --env=MYVAR=myval`,
	RunE: runStoreDeploy,
}

func runStoreDeploy(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide the function name")
	}

	storeItems, err := storeList(storeAddress)
	if err != nil {
		return err
	}

	item := storeFindFunction(args[0], storeItems)
	if item == nil {
		return fmt.Errorf("function '%s' not found", functionName)
	}

	// Add the store environment variables to the provided ones from cmd
	if item.Environment != nil {
		for k, v := range item.Environment {
			env := fmt.Sprintf("%s=%s", k, v)
			storeDeployFlags.envvarOpts = append(storeDeployFlags.envvarOpts, env)
		}
	}

	// Add the store labels to the provided ones from cmd
	if item.Labels != nil {
		for k, v := range item.Labels {
			label := fmt.Sprintf("%s=%s", k, v)
			storeDeployFlags.labelOpts = append(storeDeployFlags.labelOpts, label)
		}
	}

	// Use the network from manifest if not changed by user
	if !cmd.Flag("network").Changed {
		network = item.Network
	}

	return RunDeploy(
		args,
		item.Image,
		item.Fprocess,
		item.Name,
		storeDeployFlags,
	)
}
