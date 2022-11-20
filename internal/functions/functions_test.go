// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"bytes"
	"testing"

	"github.com/heaths/go-console"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestParamFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		stdin   string
		tty     bool
		want    string
		wantErr bool
	}{
		{
			name: "default",
			tty:  true,
			want: "world",
		},
		{
			name:  "override",
			stdin: "Earth",
			tty:   true,
			want:  "Earth",
		},
		{
			name:    "cannot prompt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		con := console.Fake(
			console.WithStdin(bytes.NewBufferString(tt.stdin+"\n")),
			console.WithStderrTTY(tt.tty),
		)

		params := make(map[string]string)
		sut := ParamFunc(con.Stdin(), con.Stderr(), con.IsStderrTTY(), params)

		t.Run(tt.name, func(t *testing.T) {
			got, err := sut("name", "world", "Who should I greet")
			if err != nil {
				if tt.wantErr {
					return
				}

				t.Fatal("unexpected error:", err)
			} else if tt.wantErr {
				t.Fatal("expected error")
			}
			assert.Equal(t, tt.want, got)

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
