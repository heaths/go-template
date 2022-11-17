// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"bytes"
	"testing"

	"github.com/heaths/go-console"
)

func TestParamFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		stdin string
		want  string
	}{
		{
			name: "default",
			want: "world",
		},
		{
			name:  "override",
			stdin: "Earth",
			want:  "Earth",
		},
	}

	for _, tt := range tests {
		con := console.Fake(
			console.WithStdin(bytes.NewBufferString(tt.stdin + "\n")),
		)

		params := make(map[string]string)
		sut := ParamFunc(con, params)

		t.Run(tt.name, func(t *testing.T) {
			got, err := sut("name", "world", "Who should I greet")
			if err != nil {
				t.Error("unexpected error:", err)
				return
			}

			if got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
				return
			}

			// Run it again and make sure the value is cached.
			got, err = sut("name", "unexpected")
			if err != nil {
				t.Fatal("unexpected error:", err)
			}

			if got != tt.want {
				t.Fatalf("want %q, got %q", tt.want, got)
			}
		})
	}
}

func TestTitle(t *testing.T) {
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

	sut := TitleFunc()
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			if got := sut(tt.value); got != tt.want {
				t.Fatalf("want %q, got %q", tt.want, got)
			}
		})
	}
}
