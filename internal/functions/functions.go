// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ParamFunc(r io.Reader, w io.Writer, isTTY bool, params map[string]string) func(string, ...any) (string, error) {
	return func(name string, args ...any) (value string, err error) {
		var ok bool
		if value, ok = params[name]; !isTTY {
			return "", fmt.Errorf("cannot prompt for parameter %q", name)

		} else if !ok {
			var valueType string
			var valid func(string) bool

			defaultValue := ""
			if len(args) > 0 {
				defaultValue, valueType, valid, err = format(args[0])
				if err != nil {
					return
				}
			}

			prompt := name
			if len(args) > 1 {
				var ok bool
				if prompt, ok = args[1].(string); !ok {
					return "", fmt.Errorf("unsupported prompt %v", args[1])
				}
				prompt = strings.TrimRightFunc(prompt, func(r rune) bool {
					return r == '?'
				})
				prompt = fmt.Sprintf("%s (%s)", prompt, name)
			}

			reader := bufio.NewReader(r)
			for {
				// Assume color support since we're on a TTY.
				fmt.Fprintf(w, "\033[32m%s? \033[90m[%s]\033[0m: ", prompt, defaultValue)

				value, err = reader.ReadString('\n')
				if err != nil {
					return
				}

				value = strings.TrimSpace(value)
				if value == "" {
					value = defaultValue
				}

				if valid(value) {
					break
				}

				fmt.Fprintf(w, "\033[31mExpected %s. Please try again.\033[0m\n", valueType)
			}

			params[name] = value
		}

		return
	}
}

func Pluralize(count int, thing string) string {
	if count == 1 {
		return fmt.Sprint(count, " ", thing)
	}

	return fmt.Sprintf("%d %ss", count, thing)
}

func PluralizeFunc(count interface{}, thing string) (string, error) {
	if i, ok := count.(int); ok {
		return Pluralize(i, thing), nil
	} else if s, ok := count.(string); ok {
		if i, err := strconv.Atoi(s); err == nil {
			return Pluralize(i, thing), nil
		} else {
			return "", err
		}
	}
	return "", fmt.Errorf("%v not a number", count)
}

func LowercaseFunc(lang language.Tag) func(string) string {
	c := cases.Lower(lang)
	return func(s string) string {
		return c.String(s)
	}
}

func TitlecaseFunc(lang language.Tag) func(string) string {
	c := cases.Title(lang)
	return func(s string) string {
		return c.String(s)
	}
}

func UppercaseFunc(lang language.Tag) func(string) string {
	c := cases.Upper(lang)
	return func(s string) string {
		return c.String(s)
	}
}

func format(v any) (value, valueType string, validate func(string) bool, err error) {
	if s, ok := v.(string); ok || v == nil {
		return s, "string", func(s string) bool { return true }, nil
	}

	if i, ok := v.(int); ok {
		return strconv.FormatInt(int64(i), 10), "integer", func(s string) bool {
			_, err := strconv.ParseInt(s, 10, 32)
			return err == nil
		}, nil
	}

	return "", "", nil, fmt.Errorf("unsupported type %v", v)
}
