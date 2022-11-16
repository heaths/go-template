// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/heaths/go-console"
	"github.com/heaths/go-template"
	"github.com/spf13/cobra"
)

func main() {
	var params []string
	cmd := &cobra.Command{
		Use:   "[flags] [root]",
		Short: "Process template files in a root directory (default is $PWD)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			pairs := make(map[string]string, len(params))
			for _, param := range params {
				tokens := strings.SplitN(param, "=", 2)
				if len(tokens) != 2 {
					return fmt.Errorf(`parameters must be specified as "name=value"`)
				}
				pairs[tokens[0]] = tokens[1]
			}

			opts := runOpts{
				con:    console.System(),
				dir:    dir,
				params: pairs,
			}

			return run(opts)
		},
	}

	cmd.Flags().StringSliceVarP(&params, "param", "p", nil, "template parameters like name=value")

	err := cmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

type runOpts struct {
	con    console.Console
	dir    string
	params map[string]string
}

func run(opts runOpts) error {
	return template.Apply(os.DirFS(opts.dir), opts.con, opts.params)
}
