// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package processor

// cspell:ignore mattn isatty
import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/heaths/go-template/internal/functions"
	"github.com/mattn/go-isatty"
	"github.com/spf13/afero"
	"golang.org/x/exp/slices"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type Processor struct {
	Stderr io.Writer // The writer on which users are prompted.
	Stdin  io.Reader // The reader from which user input is read.
	IsTTY  bool      // Whether Stderr is a terminal.

	Exclusions []string // Directories and files to exclude.

	Language *language.Tag     // The language used in some functions.
	collator *collate.Collator // The collator used to sort and search for strings.

	Log     *log.Logger // Optional logger for pertinent information.
	Verbose bool        // Whether to log verbose information.

	srcFS afero.Fs // The file system for reading templates.
	dstFS afero.Fs // The file system for writing templates.

	errors int // Number of errors logged (as warning logs).
}

func (p *Processor) Initialize() {
	if p.Stderr == nil {
		p.Stderr = os.Stderr
		p.IsTTY = isatty.IsTerminal(os.Stderr.Fd())
	}

	if p.Stdin == nil {
		p.Stdin = os.Stdin
	}

	if p.Language == nil {
		p.Language = &language.English
	}

	p.collator = collate.New(*p.Language, collate.IgnoreCase)
	p.normalizeExclusions()

	if p.srcFS == nil {
		p.srcFS = afero.NewOsFs()
	}

	if p.dstFS == nil {
		p.dstFS = p.srcFS
	}
}

func (p *Processor) Execute(root string, params map[string]string) error {
	funcs := template.FuncMap{
		"date":      functions.DateFunc,
		"lowercase": functions.LowercaseFunc(*p.Language),
		"param":     functions.ParamFunc(p.Stdin, p.Stderr, p.IsTTY, params),
		"pluralize": functions.PluralizeFunc,
		"titlecase": functions.TitlecaseFunc(*p.Language),
		"uppercase": functions.UppercaseFunc(*p.Language),
	}

	// cspell:ignore IOFS
	dir := afero.NewIOFS(p.srcFS)

	err := fs.WalkDir(dir, root, func(path string, d fs.DirEntry, err error) (_ error) {
		// TODO: Bubble up errors to channel but carry on with as many files as possible.
		if err != nil {
			p.logWarning("failed to walk %q: %v\n", path, err)
			return
		}

		switch {
		// Always ignore repos to avoid catastrophe.
		case path == ".git" || path == ".hg":
			p.logVerbose("skipping %q", path)
			return fs.SkipDir
		case p.exclude(path):
			p.logVerbose("skipping %q", path)
			if d.IsDir() {
				return fs.SkipDir
			}
			return
		case d.IsDir():
			return
		}
		p.logVerbose("processing %q", path)

		t := template.New(d.Name()).Funcs(funcs)
		t, err = t.ParseFS(dir, path)
		if err != nil {
			p.logWarning("failed to parse %q: %v\n", path, err)
			return
		}

		if !isTemplate(t) {
			p.logVerbose("skipping non-template %q", path)
			return
		}

		var file afero.File
		file, err = p.dstFS.Create(path)
		if err != nil {
			p.logWarning("failed to create output %q: %v\n", path, err)
			return
		}

		err = t.Execute(file, nil)
		if err != nil {
			p.logWarning("failed to process %q: %v\n", path, err)
			return
		}

		return
	})

	if err != nil {
		return err
	}

	if p.errors == 0 {
		return nil
	}

	return fmt.Errorf("failed to process %s", functions.Pluralize(p.errors, "template"))
}

func (p *Processor) logVerbose(format string, v ...any) {
	if p.Verbose && p.Log != nil {
		p.Log.Printf(format, v...)
	}
}

func (p *Processor) logWarning(format string, v ...any) {
	p.errors++
	if p.Log != nil {
		p.Log.Printf(format, v...)
	}
}

func (p *Processor) exclude(s string) bool {
	_, found := slices.BinarySearchFunc(p.Exclusions, s, p.collator.CompareString)
	return found
}

func (p *Processor) normalizeExclusions() {
	src := p.Exclusions
	for i, s := range src {
		s = strings.ReplaceAll(s, "\\", "/")
		s = strings.TrimRightFunc(s, func(r rune) bool {
			return r == '/'
		})

		if strings.HasPrefix(s, "/") {
			src[i] = s[1:]
		} else if strings.HasPrefix(s, "./") {
			src[i] = s[2:]
		} else {
			src[i] = s
		}
	}

	p.collator.SortStrings(src)
}

func isTemplate(t *template.Template) bool {
	for _, node := range t.Root.Nodes {
		if node.Type() != parse.NodeText {
			return true
		}
	}
	return false
}
