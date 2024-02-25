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
		if value, ok = params[name]; ok {
			// Validate the provided value if a default value was defined.
			if len(args) == 0 {
				return
			}

			var param paramValue
			param, err = fromDefaultValue(args[0])
			if err != nil {
				return
			}

			if value, ok = param.Format(value); !ok {
				return "", fmt.Errorf("invalid parameter %q value: %s; expected %s", name, value, param.Description())
			}

			return
		} else if !ok {
			if !isTTY {
				return "", fmt.Errorf("cannot prompt for parameter %q", name)
			}

			var param paramValue
			if len(args) > 0 {
				param, err = fromDefaultValue(args[0])
				if err != nil {
					return
				}
			} else {
				param = &stringValue{""}
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
				fmt.Fprintf(w, "\033[32m%s? \033[90m[%s]\033[0m: ", prompt, param.Display())

				value, err = reader.ReadString('\n')
				if err != nil {
					return
				}

				value = strings.TrimSpace(value)
				if value == "" {
					value = param.String()
					break
				}

				if value, ok = param.Format(value); ok {
					break
				}

				fmt.Fprintf(w, "\033[31mExpected %s. Please try again.\033[0m\n", param.Description())
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

func Replace(from, to, source string) string {
	return strings.Replace(source, from, to, -1)
}

func DeleteFunc(delete *bool) func() string {
	return func() string {
		*delete = true
		return ""
	}
}
