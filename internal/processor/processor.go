// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package processor

// cspell:ignore mattn isatty
import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"text/template"

	"github.com/heaths/go-template/internal/functions"
	"github.com/mattn/go-isatty"
	"golang.org/x/text/language"
)

type Processor struct {
	Stderr io.Writer // The writer on which users are prompted.
	Stdin  io.Reader // The reader from which user input is read.
	IsTTY  bool      // Whether Stderr is a terminal.

	Language *language.Tag // The language used in some functions.

	OutDir fs.FS // Output directory where files will be written.

	Log     *log.Logger // Optional logger for pertinent information.
	Verbose bool        // Whether to log verbose information.

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

	if p.Log == nil {
		p.Log = log.New(p.Stderr, "", log.Ltime)
	}
}

func (p Processor) Execute(root fs.FS, params map[string]string) error {
	funcs := template.FuncMap{
		"param":     functions.ParamFunc(p.Stdin, p.Stderr, p.IsTTY, params),
		"titlecase": functions.TitleFunc(*p.Language),
	}

	// TODO: Parallelize work into sync.WaitGroup or errgroup.Group for CPU count.
	err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) (_ error) {
		// TODO: Bubble up errors to channel but carry on with as many files as possible.
		if err != nil {
			p.warn("failed to walk %q: %v\n", path, err)
			return
		}

		switch {
		// TODO: Take options for exclusions including some defaults, or func for caller to to validate.
		case d.Name() == ".git":
			p.info("skipping %q", path)
			return
		case d.IsDir():
			return
		}

		t := template.New(d.Name()).Funcs(funcs)
		t, err = t.ParseFS(root, path)
		if err != nil {
			p.warn("failed to parse %q: %v\n", path, err)
			return
		}

		// TODO: Write to "file~" and move to "file" if successful.
		stdout := &bytes.Buffer{}
		err = t.Execute(stdout, nil)
		if err != nil {
			p.warn("failed to process %q: %v\n", path, err)
			return
		}

		_, err = io.Copy(os.Stdout, stdout)
		if err != nil {
			p.warn("failed to write %q: %v\n", path, err)
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

	return fmt.Errorf("failed to process %s", pluralize(p.errors, "file"))
}

func (p Processor) info(format string, v ...any) {
	if p.Verbose {
		p.Log.Printf(format, v...)
	}
}

func (p *Processor) warn(format string, v ...any) {
	p.errors++
	p.Log.Printf(format, v...)
}

func pluralize(count int, thing string) string {
	if count == 1 {
		return fmt.Sprint(count, thing)
	}

	return fmt.Sprintf("%d %ss", count, thing)
}
