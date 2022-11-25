// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package template

import (
	"io"
	"log"

	"github.com/heaths/go-template/internal/processor"
	"golang.org/x/text/language"
)

// ApplyOption applies options to the processor.
type ApplyOption func(*processor.Processor)

// Apply applies parameters to all templates with the given root directory.
func Apply(root string, params map[string]string, options ...ApplyOption) error {
	proc := new(processor.Processor)
	for _, opt := range options {
		opt(proc)
	}
	proc.Initialize()

	return proc.Execute(root, params)
}

// WithOutput specifies the output Writer and whether it represents a TTY.
// By default this is os.Stderr. isTTY depends on whether os.Stderr
// is redirected.
func WithOutput(w io.Writer, isTTY bool) ApplyOption {
	return func(p *processor.Processor) {
		p.Stderr = w
		p.IsTTY = isTTY
	}
}

// WithInput specifies the input Reader. By default this is os.Stdin.
func WithInput(r io.Reader) ApplyOption {
	return func(p *processor.Processor) {
		p.Stdin = r
	}
}

// WithDelims specifies alternate delimiters to open and close template expressions.
// The defaults are "{{" and "}}". Both or neither must be non-empty or this function panics.
func WithDelims(left, right string) ApplyOption {
	if left != right && (left == "" || right == "") {
		panic("both or neither left and right must be non-empty")
	}

	return func(p *processor.Processor) {
		p.LeftDelim = left
		p.RightDelim = right
	}
}

// WithExclusions specifies excluded directories and files. These paths should be
// relative to the root directory passed to Apply. The prefixes "./" and "/" are
// automatically removed. Comparisons are case-insensitive.
func WithExclusions(exclusions []string) ApplyOption {
	return func(p *processor.Processor) {
		p.Exclusions = exclusions
	}
}

// WithLanguage specifies the language for any template function that needs it.
// The default is language.English.
func WithLanguage(language language.Tag) ApplyOption {
	return func(p *processor.Processor) {
		p.Language = &language
	}
}

// WithLogger specifies the logger to write to and whether to log verbose output.
// No logging is performed by default.
func WithLogger(log *log.Logger, verbose bool) ApplyOption {
	return func(p *processor.Processor) {
		p.Log = log
		p.Verbose = verbose
	}
}
