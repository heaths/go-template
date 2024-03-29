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
			name:         "re-prompt (empty default)",
			defaultValue: "",
			stdin:        "world",
			tty:          true,
			want:         "world",
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
		{
			name:         "unsupported param (no TTY)",
			defaultValue: time.Now,
			param:        "2022-11-25",
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

func TestReplace(t *testing.T) {
	t.Parallel()

	sut := Replace("-", "_", "my-crate")
	assert.Equal(t, "my_crate", sut)
}

func TestDelete(t *testing.T) {
	t.Parallel()

	current := "current"
	var delete bool
	var values []string
	sut := DeleteFunc(&current, &delete, &values)

	assert.Equal(t, "", sut())
	assert.True(t, delete)
	assert.Equal(t, []string{"current"}, values)

	assert.Equal(t, "", sut("foo", "bar"))
	assert.True(t, delete)
	assert.Equal(t, []string{"current", "foo", "bar"}, values)

	assert.Equal(t, "", sut("baz"))
	assert.True(t, delete)
	assert.Equal(t, []string{"current", "foo", "bar", "baz"}, values)
}
