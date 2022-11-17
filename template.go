// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package template

import (
	"io"
	"io/fs"
	"log"

	"github.com/heaths/go-template/internal/processor"
	"golang.org/x/text/language"
)

// Options to pass to Apply.
type ApplyOption func(*processor.Processor)

// Applies parameters to all templates with the given root directory.
func Apply(root fs.FS, params map[string]string, options ...ApplyOption) error {
	proc := new(processor.Processor)
	for _, opt := range options {
		opt(proc)
	}
	proc.Initialize()

	return proc.Execute(root, params)
}

// Specify the output Writer and whether it represents a TTY.
// By default this is os.Stderr. isTTY depends on whether os.Stderr
// is redirected.
func WithOutput(w io.Writer, isTTY bool) ApplyOption {
	return func(p *processor.Processor) {
		p.Stderr = w
		p.IsTTY = isTTY
	}
}

// Specify the input Reader. By default this is os.Stdin.
func WithInput(r io.Reader) ApplyOption {
	return func(p *processor.Processor) {
		p.Stdin = r
	}
}

// Specify the language for any template function that needs it.
// By default this is language.English.
func WithLanguage(language language.Tag) ApplyOption {
	return func(p *processor.Processor) {
		*p.Language = language
	}
}

// Specify the logger to write to and whether to log verbose output.
// By default this is log.Default() without verbose logging.
func WithLogger(log *log.Logger, verbose bool) ApplyOption {
	return func(p *processor.Processor) {
		p.Log = log
		p.Verbose = verbose
	}
}
