// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package template

import (
	"bytes"
	"io"
	_fs "io/fs"
	"log"
	"text/template"

	"github.com/heaths/go-console"
	"github.com/heaths/go-template/internal/functions"
)

func Apply(fs _fs.FS, con console.Console, params map[string]string) error {
	// TODO: Pass (optional?) io.Writer for stderr and io.Reader for stdin. Use console only for apps, tests.
	funcs := template.FuncMap{
		"param":     functions.ParamFunc(con, params),
		"titlecase": functions.TitleFunc(),
	}

	err := _fs.WalkDir(fs, ".", func(path string, d _fs.DirEntry, err error) (_ error) {
		// TODO: Bubble up errors to channel but carry on with as many files as possible.
		if err != nil {
			log.Printf("failed to walk %q: %v\n", path, err)
			return
		}

		switch {
		// TODO: Take options for exclusions including some defaults, or func for caller to to validate.
		case d.Name() == ".git":
			log.Printf("skipping %q", path)
			return
		case d.IsDir():
			return
		}

		t := template.New(d.Name()).Funcs(funcs)
		t, err = t.ParseFS(fs, path)
		if err != nil {
			log.Printf("failed to parse %q: %v\n", path, err)
			return
		}

		// TODO: Write to "file~" and move to "file" if successful.
		stdout := &bytes.Buffer{}
		err = t.Execute(stdout, nil)
		if err != nil {
			log.Printf("failed to process %q: %v\n", path, err)
			return
		}

		_, err = io.Copy(con.Stdout(), stdout)
		if err != nil {
			log.Printf("failed to write %q: %v\n", path, err)
			return
		}

		return
	})

	return err
}
