// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ParamFunc(r io.Reader, w io.Writer, isTTY bool, params map[string]string) func(string, ...string) (string, error) {
	return func(name string, args ...string) (value string, err error) {
		var ok bool
		if value, ok = params[name]; !isTTY {
			return "", fmt.Errorf("cannot prompt for parameter %q", name)

		} else if !ok {
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

			fmt.Fprintf(w, "\033[32m%s? \033[90m[%s]\033[0m: ", prompt, defaultValue)

			reader := bufio.NewReader(r)
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

func TitleFunc(lang language.Tag) func(string) string {
	c := cases.Title(lang)
	return func(s string) string {
		return c.String(s)
	}
}
