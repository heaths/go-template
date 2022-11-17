// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"bufio"
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/heaths/go-console"
	"github.com/heaths/go-console/pkg/colorscheme"
)

func ParamFunc(con console.Console, params map[string]string) func(string, ...string) (string, error) {
	cs := colorscheme.New(colorscheme.WithTTY(con.IsStderrTTY))
	return func(name string, args ...string) (value string, err error) {
		var ok bool
		if value, ok = params[name]; !ok {
			defaultValue := ""
			if len(args) > 0 {
				defaultValue = args[0]
			}

			prompt := name
			if len(args) > 1 {
				prompt = strings.TrimRightFunc(args[1], func(r rune) bool {
					return r == '?'
				})
				prompt = fmt.Sprintf("%s (%s)", prompt, name)
			}

			fmt.Fprintf(con.Stderr(), cs.Green("%s? ")+cs.LightBlack("[%s] "), prompt, defaultValue)

			reader := bufio.NewReader(con.Stdin())
			value, err = reader.ReadString('\n')
			if err != nil {
				return
			}

			value = strings.TrimSpace(value)
			if value == "" {
				value = defaultValue
			}

			params[name] = value
		}

		return
	}
}

func TitleFunc() func(string) string {
	// TODO: Allow passing in language tag.
	c := cases.Title(language.English)
	return func(s string) string {
		return c.String(s)
	}
}
