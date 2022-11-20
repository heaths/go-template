# Apply Templates

Process all template files recursively in a directory. This is useful for
project templates, for example, after creating a new repository from a template
repository.

## Example

All files in the specified directory will be processed as templates. For now,
this project assumes all files are UTF-8 encoded. A _.git_ directory or file
(worktree) will be skipped.

For every `{{param}}` found in a template, the user will be prompted to answer
unless the named parameter was already in the parameter cache you pass. This
cache can be pre-populated as well.

```golang
import (
    "log"

    "github.com/heaths/go-template"
)

func Example() {
    params := make(map[string]string)
    err := template.Apply("testdata", params,
        template.WithLogger(log.Default(), true),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## Templates

Templates are processed using [`text/template`](https://pkg.go.dev/text/template).
Template files contain a mix of text and actions surrounded by `{{` and `}}` e.g.,

```markdown
# {{param "name" "" "What is the project name?" | titlecase}}

This is an example.
```

### Functions

In addition to [built-in](https://pkg.go.dev/text/template#hdr-Functions) functions,
the following functions are also available:

* `param <name> [<default> [<prompt>]]`\
  Replace with a parameter named `<name>`, or prompt using an optional `<default>`
  with an optional `<prompt>`. If a `<prompt>` is not specified, the required
  `<name>` is used.
* `pluralize <count> <thing>`\
  Append an "s" to `<thing>` if `<count>` is not equal to 1. `<count>` can be
  either an `int` or a `string` representing an `int` e.g., "1".
* `lowercase <string>`\
  Change the case of `<string>` to all lowercase characters.
* `titlecase <string>`\
  Change the case of `<string>` to Title Case characters.
* `uppercase <string>`\
  Change the case of `<string>` to UPPERCASE characters.

## License

Licensed under the [MIT](LICENSE.txt) license.
