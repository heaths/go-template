// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package template_test

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/heaths/go-console"
	"github.com/heaths/go-template"
)

func ExampleApply() {
	testdata := os.DirFS("./testdata")
	con := console.Fake(
		console.WithStderrTTY(true),
		console.WithStdin(bytes.NewBufferString("template\n")),
	)
	params := make(map[string]string)

	err := template.Apply(testdata, params,
		template.WithOutput(con.Stderr(), con.IsStderrTTY()),
		template.WithInput(con.Stdin()),
	)
	if err != nil {
		log.Fatalln("failed to apply templates:", err)
	}

	stdout, _, _ := con.Buffers()
	fmt.Print(stdout.String())

	// Output:
	// # Template
	//
	// Project "Template" is an example.
}
