// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"bytes"
	"testing"
	"time"

	"github.com/heaths/go-console"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestParamFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		defaultValue interface{}
		stdin        string
		tty          bool
		param        string
		want         string
		wantErr      bool
	}{
		{
			name:         "default",
			defaultValue: "world",
			tty:          true,
			want:         "world",
		},
		{
			name:         "override",
			defaultValue: "world",
			stdin:        "Earth",
			tty:          true,
			want:         "Earth",
		},
		{
			name:         "cannot prompt",
			defaultValue: "world",
			wantErr:      true,
		},
		{
			name:         "re-prompt (default)",
			defaultValue: 2022,
			stdin:        "world\n",
			tty:          true,
			want:         "2022",
		},
		{
			name:  "re-prompt (no default)",
			stdin: "world",
			tty:   true,
			want:  "world",
		},
		{
			name:         "re-prompt",
			defaultValue: 2022,
			stdin:        "world\n2023",
			tty:          true,
			want:         "2023",
		},
		{
			name:         "boolean (default true)",
			defaultValue: true,
			tty:          true,
			want:         "true",
		},
		{
			name:         "boolean (no)",
			defaultValue: true,
			stdin:        "no",
			tty:          true,
			want:         "",
		},
		{
			name:         "integer param (no TTY)",
			defaultValue: 1,
			param:        "2",
			want:         "2",
		},
		{
			name:         "invalid integer param (no TTY)",
			defaultValue: 1,
			param:        "invalid",
			wantErr:      true,
		},
		{
			name:         "unsupported",
			defaultValue: time.Now,
			tty:          true,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		con := console.Fake(
			console.WithStdin(bytes.NewBufferString(tt.stdin+"\n")),
			console.WithStderrTTY(tt.tty),
		)
		_, stderr, _ := con.Buffers()

		params := make(map[string]string)
		if tt.param != "" {
			params["name"] = tt.param
		}
		sut := ParamFunc(con.Stdin(), con.Stderr(), con.IsStderrTTY(), params)

		t.Run(tt.name, func(t *testing.T) {
			got, err := sut("name", tt.defaultValue, "What should I prompt?")
			if err != nil {
				if tt.wantErr {
					return
				}

				t.Fatal("unexpected error:", err)
			} else if tt.wantErr {
				t.Fatal("expected error")
			}

			if tt.tty {
				assert.Contains(t, stderr.String(), "What should I prompt (")
			}

			if !assert.Equal(t, tt.want, got) {
				return
			}

			// Run it again and make sure the value is cached.
			got, err = sut("name", "unexpected")
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPluralize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		count int
		want  string
	}{
		{
			name: "zero",
			want: "0 things",
		},
		{
			name:  "singular",
			count: 1,
			want:  "1 thing",
		},
		{
			name:  "plural",
			count: 2,
			want:  "2 things",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pluralize(tt.count, "thing")
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPluralizeFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		count   interface{}
		want    string
		wantErr bool
	}{
		{
			name:  "plural int",
			count: 2,
			want:  "2 things",
		},
		{
			name:  "singular string",
			count: "1",
			want:  "1 thing",
		},
		{
			name:    "nan string",
			count:   "nan",
			wantErr: true,
		},
		{
			name:    "bool",
			count:   true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PluralizeFunc(tt.count, "thing")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLowercase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value string
		want  string
	}{
		{
			value: "lOrD oF tHe RiNgs",
			want:  "lord of the rings",
		},
		{
			value: "the hobbit",
			want:  "the hobbit",
		},
	}

	sut := LowercaseFunc(language.English)
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got := sut(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTitlecase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value string
		want  string
	}{
		{
			value: "lOrD oF tHe RiNgs",
			// Expected output is wrong for English or AmericanEnglish.
			want: "Lord Of The Rings",
		},
		{
			value: "the hobbit",
			want:  "The Hobbit",
		},
	}

	sut := TitlecaseFunc(language.English)
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got := sut(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUppercase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value string
		want  string
	}{
		{
			value: "lOrD oF tHe RiNgs",
			want:  "LORD OF THE RINGS",
		},
		{
			value: "the hobbit",
			want:  "THE HOBBIT",
		},
	}

	sut := UppercaseFunc(language.English)
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got := sut(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		value        interface{}
		other        string
		isOtherValid bool
		wantValue    string
		wantType     string
		wantDisp     string
		wantErr      bool
	}{
		{
			name:         "string",
			value:        "value",
			other:        "other",
			wantValue:    "value",
			wantType:     "a string",
			isOtherValid: true,
		},
		{
			name:         "int",
			value:        1,
			other:        "2",
			wantValue:    "1",
			wantType:     "an integer",
			isOtherValid: true,
		},
		{
			name:      "int (invalid)",
			value:     1,
			other:     "other",
			wantValue: "1",
			wantType:  "an integer",
		},
		{
			name:         "boolean (true)",
			value:        true,
			other:        "y",
			wantValue:    "true",
			wantType:     "yes (Y) or no (N)",
			wantDisp:     "Y/n",
			isOtherValid: true,
		},
		{
			name:         "boolean (false)",
			value:        false,
			other:        "n",
			wantValue:    "",
			wantType:     "yes (Y) or no (N)",
			wantDisp:     "y/N",
			isOtherValid: true,
		},
		{
			name:      "boolean (invalid)",
			value:     false,
			other:     "invalid",
			wantValue: "",
			wantType:  "yes (Y) or no (N)",
			wantDisp:  "y/N",
		},
		{
			name:    "time (unsupported)",
			value:   time.Now,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, typeDescription, display, format, err := convert(tt.value)
			if tt.wantErr {
				assert.Errorf(t, err, "unsupported type")
				return
			} else if !assert.NoError(t, err) {
				return
			}

			_, ok := format(value)
			assert.True(t, ok)
			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantType, typeDescription)

			if tt.wantDisp == "" {
				tt.wantDisp = tt.wantValue
			}

			_, ok = format(tt.other)
			assert.Equal(t, tt.isOtherValid, ok)
			assert.Equal(t, tt.wantDisp, display())
		})
	}
}
