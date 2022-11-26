// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"fmt"
	"strconv"
	"strings"
)

type paramValue interface {
	fmt.Stringer

	Description() string
	Display() string
	Format(string) (string, bool)
}

func fromDefaultValue(v any) (paramValue, error) {
	switch v := v.(type) {
	case string:
		return &stringValue{v}, nil
	case int:
		return &intValue{v}, nil
	case bool:
		return &boolValue{v}, nil
	default:
		return nil, fmt.Errorf("unsupported type %v", v)
	}
}

type stringValue struct {
	defaultValue string
}

func (v stringValue) Description() string {
	return "a string"
}

func (v stringValue) Display() string {
	return v.String()
}

func (v stringValue) Format(s string) (string, bool) {
	return s, true
}

func (v stringValue) String() string {
	return v.defaultValue
}

type intValue struct {
	defaultValue int
}

func (v intValue) Description() string {
	return "an integer"
}

func (v intValue) Display() string {
	return v.String()
}

func (v intValue) Format(s string) (string, bool) {
	_, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return "", false
	}
	return s, true
}

func (v intValue) String() string {
	return strconv.FormatInt(int64(v.defaultValue), 10)
}

type boolValue struct {
	defaultValue bool
}

func (v boolValue) Description() string {
	return "yes (Y) or no (N)"
}

func (v boolValue) Display() string {
	if v.defaultValue {
		return "Y/n"
	}
	return "y/N"
}

func (v boolValue) Format(s string) (string, bool) {
	if s == "" {
		return v.String(), true
	}
	if strings.EqualFold(s, "y") || strings.EqualFold(s, "yes") || strings.EqualFold(s, "true") {
		return "true", true
	}
	if strings.EqualFold(s, "n") || strings.EqualFold(s, "no") || strings.EqualFold(s, "false") {
		// text/template's `if` treats zero values as false.
		return "", true
	}
	return "", false
}

func (v boolValue) String() string {
	if v.defaultValue {
		return "true"
	}
	// text/template's `if` treats zero values as false.
	return ""
}
