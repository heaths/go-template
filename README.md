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
    "os"

    "github.com/heaths/go-console"
    "github.com/heaths/go-template"
)

func Example() {
    dir := os.DirFS(".")
    con := console.System()
    params := make(map[string]string)

    err := template.Apply(dir, con, params)
    if err != nil {
        log.Fatalln("failed to process templates:", err)
    }
}
```

## License

Licensed under the [MIT](LICENSE.txt) license.
