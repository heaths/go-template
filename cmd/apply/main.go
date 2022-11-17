// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"log"

	"github.com/heaths/go-template"
	"github.com/spf13/cobra"
)

func main() {
	var params map[string]string
	verbose := false
	cmd := &cobra.Command{
		Use:   "[flags] [root]",
		Short: "Process template files in a root directory (default is $PWD)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) > 0 {
				root = args[0]
			}

			if params == nil {
				params = make(map[string]string)
			}

			return template.Apply(
				root,
				params,
				template.WithLogger(log.Default(), verbose),
			)
		},
	}

	cmd.Flags().StringToStringVarP(&params, "param", "p", nil, "template parameters like name=value")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "log verbose output")

	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
