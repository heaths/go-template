// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package template

import (
	_fs "io/fs"

	"github.com/heaths/go-console"
	"github.com/heaths/go-template/internal/processor"
)

func Apply(fs _fs.FS, con console.Console, params map[string]string) error {
	proc := processor.Processor{
		Stderr: con.Stderr(),
		Stdin:  con.Stdin(),
		IsTTY:  con.IsStderrTTY(),
	}
	proc.Initialize()

	return proc.Execute(fs, params)
}
