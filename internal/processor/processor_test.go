// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package processor

import (
	"bytes"
	"io"
	"testing"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	"github.com/heaths/go-console"
	"github.com/spf13/afero"
)

const (
	a = `# {{param "name" "" "What is the project name?" | titlecase}}

Project "{{param "name" | titlecase}}" is an example.
`
)

func TestProcessor_Execute(t *testing.T) {
	t.Parallel()

	con := console.Fake(
		console.WithStdin(bytes.NewBufferString("template\n")),
		console.WithStderrTTY(true),
	)

	srcFS := afero.NewMemMapFs()
	srcFS.Mkdir("testdata", 0755)
	afero.WriteFile(srcFS, "testdata/a.md", []byte(a), 0644)

	dstFS := afero.NewMemMapFs()

	proc := Processor{
		Stderr: con.Stderr(),
		Stdin:  con.Stdin(),
		IsTTY:  con.IsStderrTTY(),

		srcFS: srcFS,
		dstFS: afero.NewCopyOnWriteFs(srcFS, dstFS),
	}
	proc.Initialize()

	params := make(map[string]string)
	err := proc.Execute("testdata", params)
	if err != nil {
		t.Fatalf("failed to process template: %v", err)
	}

	const path = "testdata/a.md"
	file, err := dstFS.Open(path)
	if err != nil {
		t.Fatalf("failed to open %q: %v", path, err)
	}

	got, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to read %q: %v", path, err)
	}

	want := heredoc.Doc(`
		# Template

		Project "Template" is an example.
		`)

	if string(got) != want {
		t.Fatalf("want %q, got %q", want, string(got))
	}
}

func TestIsTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		template string
		want     bool
	}{
		{
			name:     "template",
			template: `Hello, {{"world"}}!`,
			want:     true,
		},
		{
			name:     "not template",
			template: "Hello, world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut, err := template.New(tt.name).Parse(tt.template)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := isTemplate(sut); got != tt.want {
				t.Fatalf("want %t, got %t", tt.want, got)
			}
		})
	}
}
