// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package processor

import (
	"bytes"
	"io"
	"strconv"
	"testing"
	"text/template"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/heaths/go-console"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

const (
	content_a = `# {{param "name" "" "What is the project name?" | titlecase}}

Project "{{param "name" | titlecase}}" is an example of template repository {{param "github.owner"}}/{{param "github.repo"}}.

Copyright {{date.Local.Year}} {{param "git.name"}} under the [MIT](LICENSE.txt) license.
`

	content_a_alt = `# <%param "name" "" "What is the project name?" | titlecase%>

Project "<%param "name" | titlecase%>" is an example of template repository <%param "github.owner"%>/<%param "github.repo"%>.

Copyright <%date.Local.Year%> <%param "git.name"%> under the [MIT](LICENSE.txt) license.
`
)

// cspell:ignore Docf
func TestProcessor_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		leftDelim  string
		rightDelim string
		content    string
	}{
		{
			name:    "defaults",
			content: content_a,
		},
		{
			name:       "alternate delims",
			leftDelim:  "<%",
			rightDelim: "%>",
			content:    content_a_alt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			assert.NoError(t, afero.WriteFile(srcFS, "testdata/a.md", []byte(tt.content), 0644))
			assert.NoError(t, afero.WriteFile(srcFS, "testdata/b.md", []byte("not a template"), 0644))

			dstFS := afero.NewMemMapFs()

			proc := Processor{
				Stderr: con.Stderr(),
				Stdin:  con.Stdin(),
				IsTTY:  con.IsStderrTTY(),

				LeftDelim:  tt.leftDelim,
				RightDelim: tt.rightDelim,
				Exclusions: []string{"Build/"},

				srcFS: srcFS,
				dstFS: afero.NewCopyOnWriteFs(srcFS, dstFS),
			}
			proc.Initialize()

			params := map[string]string{
				"git.name":     "Heath Stewart",
				"github.owner": "heaths",
				"github.repo":  "template-golang",
			}
			err = proc.Execute(".", params)
			assert.NoError(t, err, "failed to process template")

			_, err = dstFS.Stat(".git")
			assert.Error(t, err)

			_, err = dstFS.Stat("build")
			assert.Error(t, err)

			_, err = dstFS.Stat("testdata/b.md")
			assert.Error(t, err)

			const path = "testdata/a.md"
			file, err := dstFS.Open(path)
			if !assert.NoError(t, err, "failed to open %q", path) {
				return
			}

			got, err := io.ReadAll(file)
			assert.NoError(t, err, "failed to read %q", path)

			// There's a small but acceptable window where the year could be different due to TZ offset.
			want := heredoc.Docf(`
				# Template

				Project "Template" is an example of template repository heaths/template-golang.

				Copyright %s Heath Stewart under the [MIT](LICENSE.txt) license.
				`, strconv.FormatInt(int64(time.Now().UTC().Year()), 10))

			assert.Equal(t, want, string(got))
		})
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
			assert.NoError(t, err)

			got := isTemplate(sut)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessor_exclude(t *testing.T) {
	t.Parallel()

	src := []string{
		"/testdata/B",
		"./testdata/a",
		"build\\c",
		"Dist/",
	}

	p := Processor{
		Exclusions: src,
		collator:   collate.New(language.English, collate.IgnoreCase),
	}

	p.normalizeExclusions()
	assert.True(t, p.exclude("testdata/b"))
}

func TestProcessor_normalizeExclusions(t *testing.T) {
	t.Parallel()

	src := []string{
		"/testdata/B",
		"./testdata/a",
		"build\\c",
		"Dist/",
	}

	dst := []string{
		"build/c",
		"Dist",
		"testdata/a",
		"testdata/B",
	}

	p := Processor{
		Exclusions: src,
		collator:   collate.New(language.English, collate.IgnoreCase),
	}

	p.normalizeExclusions()
	assert.Equal(t, dst, src)
}
