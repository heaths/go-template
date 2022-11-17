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
	"text/template"
	"text/template/parse"

	"github.com/heaths/go-template/internal/functions"
	"github.com/mattn/go-isatty"
	"github.com/spf13/afero"
	"golang.org/x/text/language"
)

type Processor struct {
	Stderr io.Writer // The writer on which users are prompted.
	Stdin  io.Reader // The reader from which user input is read.
	IsTTY  bool      // Whether Stderr is a terminal.

	Language *language.Tag // The language used in some functions.

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

	if p.srcFS == nil {
		p.srcFS = afero.NewOsFs()
	}

	if p.dstFS == nil {
		p.dstFS = p.srcFS
	}
}

func (p *Processor) Execute(root string, params map[string]string) error {
	funcs := template.FuncMap{
		"param":     functions.ParamFunc(p.Stdin, p.Stderr, p.IsTTY, params),
		"lowercase": functions.LowercaseFunc(*p.Language),
		"titlecase": functions.TitlecaseFunc(*p.Language),
		"uppercase": functions.UppercaseFunc(*p.Language),
		// TODO: accept string from `param` for count.
		// "pluralize": functions.Pluralize,
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
		// TODO: Take options for exclusions including some defaults, or func for caller to to validate.
		case d.Name() == ".git":
			p.logVerbose("skipping %q", path)
			return fs.SkipDir
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

func isTemplate(t *template.Template) bool {
	for _, node := range t.Root.Nodes {
		if node.Type() != parse.NodeText {
			return true
		}
	}
	return false
}
