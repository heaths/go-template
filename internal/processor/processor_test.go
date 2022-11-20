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
	"github.com/stretchr/testify/assert"
)

const (
	a = `# {{param "name" "" "What is the project name?" | titlecase}}

Project "{{param "name" | titlecase}}" is an example of template repository {{param "github.owner"}}/{{param "github.repo"}}.
`
)

func TestProcessor_Execute(t *testing.T) {
	t.Parallel()

	var err error
	con := console.Fake(
		console.WithStdin(bytes.NewBufferString("template\n")),
		console.WithStderrTTY(true),
	)

	srcFS := afero.NewMemMapFs()
	assert.NoError(t, srcFS.Mkdir(".git", 0755))
	assert.NoError(t, afero.WriteFile(srcFS, ".git/index", []byte("Head: main"), 0644))
	assert.NoError(t, srcFS.Mkdir("build", 0755))
	assert.NoError(t, afero.WriteFile(srcFS, "build/dat", []byte{00, 01, 02, 03}, 0644))
	assert.NoError(t, srcFS.Mkdir("testdata", 0755))
	assert.NoError(t, afero.WriteFile(srcFS, "testdata/a.md", []byte(a), 0644))

	dstFS := afero.NewMemMapFs()

	proc := Processor{
		Stderr: con.Stderr(),
		Stdin:  con.Stdin(),
		IsTTY:  con.IsStderrTTY(),

		Exclusions: []string{"build/"},

		srcFS: srcFS,
		dstFS: afero.NewCopyOnWriteFs(srcFS, dstFS),
	}
	proc.Initialize()

	params := map[string]string{
		"github.owner": "heaths",
		"github.repo":  "template-golang",
	}
	err = proc.Execute("testdata", params)
	assert.NoError(t, err, "failed to process template")

	_, err = dstFS.Stat(".git")
	assert.Error(t, err)

	_, err = dstFS.Stat("build")
	assert.Error(t, err)

	const path = "testdata/a.md"
	file, err := dstFS.Open(path)
	assert.NoError(t, err, "failed to open %q", path)

	got, err := io.ReadAll(file)
	assert.NoError(t, err, "failed to read %q", path)

	want := heredoc.Doc(`
		# Template

		Project "Template" is an example of template repository heaths/template-golang.
		`)

	assert.Equal(t, want, string(got))
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
			assert.NoError(t, err)

			got := isTemplate(sut)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNormalizeExclusions(t *testing.T) {
	src := []string{
		"/testdata/b",
		"./testdata/a",
		"build\\c",
		"dist/",
	}

	// Not currently sorted. See comment in normalizeExclusions.
	dst := []string{
		"testdata/b",
		"testdata/a",
		"build/c",
		"dist",
	}

	normalizeExclusions(src)
	assert.Equal(t, dst, src)
}
