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

			var typeDescription string
			var format func(string) (string, bool)
			_, typeDescription, _, format, err = convert(args[0])

			if value, ok = format(value); !ok {
				return "", fmt.Errorf("invalid parameter %q value: %s; expected %s", name, value, typeDescription)
			}

			return
		} else if !ok {
			if !isTTY {
				return "", fmt.Errorf("cannot prompt for parameter %q", name)
			}

			var typeDescription string
			var defaultDisplay func() string
			var format func(string) (string, bool)

			defaultValue := ""
			if len(args) > 0 {
				defaultValue, typeDescription, defaultDisplay, format, err = convert(args[0])
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
				fmt.Fprintf(w, "\033[32m%s? \033[90m[%s]\033[0m: ", prompt, defaultDisplay())

				value, err = reader.ReadString('\n')
				if err != nil {
					return
				}

				value = strings.TrimSpace(value)
				if value == "" {
					value = defaultValue
					break
				}

				if value, ok = format(value); ok {
					break
				}

				fmt.Fprintf(w, "\033[31mExpected %s. Please try again.\033[0m\n", typeDescription)
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

func convert(v any) (defaultValue, typeDescription string, defaultDisplay func() string, format func(string) (string, bool), err error) {
	if s, ok := v.(string); ok || v == nil {
		return s,
			"a string",
			func() string { return s },
			func(s string) (string, bool) { return s, true },
			nil
	}

	if i, ok := v.(int); ok {
		id := func() string { return strconv.FormatInt(int64(i), 10) }
		return id(),
			"an integer",
			id,
			func(s string) (string, bool) {
				_, err := strconv.ParseInt(s, 10, 32)
				return s, err == nil
			},
			nil
	}

	if b, ok := v.(bool); ok {
		// text/template's `if` treats zero values as false.
		id := func() string {
			if b {
				return "true"
			}
			return ""
		}
		return id(),
			"yes (Y) or no (N)",
			func() string {
				if b {
					return "Y/n"
				}
				return "y/N"
			},
			func(s string) (string, bool) {
				if s == "" {
					return id(), true
				}
				if strings.EqualFold(s, "y") ||
					strings.EqualFold(s, "yes") ||
					strings.EqualFold(s, "true") {
					return "true", true
				}
				if strings.EqualFold(s, "n") ||
					strings.EqualFold(s, "no") ||
					strings.EqualFold(s, "false") {
					return "", true
				}
				return "", false
			},
			nil
	}

	return "", "", nil, nil, fmt.Errorf("unsupported type %v", v)
}
