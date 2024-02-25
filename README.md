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
  `<name>` is used. The type of `<default>` dictates valid input. Only `string`
  and `int` are supported at this time.
* `pluralize <count> <thing>`\
  Append an "s" to `<thing>` if `<count>` is not equal to 1. `<count>` can be
  either an `int` or a `string` representing an `int` e.g., "1".
* `lowercase <string>`\
  Change the case of `<string>` to all lowercase characters.
* `titlecase <string>`\
  Change the case of `<string>` to Title Case characters.
* `uppercase <string>`\
  Change the case of `<string>` to UPPERCASE characters.
* `replace <from> <to> <source>`\
  Replaces all occurrences of `<from>` to `<to>` in the `<source>` string.
* `date`\
  Returns the current UTC date-time.
* `date.Format <layout>`\
  Formats the date-time according to [`time.Format`](https://pkg.go.dev/time#Time.Format).
* `date.Local`\
  Returns the current local date-time. You can call other `date` functions
  on the returned value e.g., `date.Local.Year`.
* `date.Year`\
  Returns the current UTC year.
* `true`\
  Returns `true`. Useful as a default value to accept yes/Y or no/N answers.
* `false`\
  Returns `false`. Useful as a default value to accept yes/Y or no/N answers.
* `deleteFile`\
  Deletes the current file.

Note that `date` functions including `Format`, `Local`, and `Year` are function calls
and need to be closed in parenthesis if you want to pipe to another function like `printf`:

```text
{{param "copyright" ((date.Year) | printf "Copyright %d") "What is the copyright year?"}}
```

To require an integer when prompting for the `copyright` parameter,
a better example is to pass the `int` that `date.Year` returns:

```text
Copyright {{param "copyright" (date.Year) "What is the copyright year?"}}
```

## License

Licensed under the [MIT](LICENSE.txt) license.
