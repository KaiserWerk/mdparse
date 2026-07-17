# mdparse

`mdparse` is a minimalistic markdown parser written in Go. Use it to get separate frontmatter and markdown content from a markdown file.
The markdown is usually needed as-is, but the frontmatter needs to be unmarshalled. This is that `mdparse` offers.

There is only two functions:

```go
func ParseFile(file string, markdown *string, v any) error
func Parse(input string, markdown *string, v any) error
```
Internally, `ParseFile` just calls `Parse` after reading the file.

## Usage

It can be used like this:

```go
type frontmatter struct {
    Title string `yaml:"title"`
    Date  string `yaml:"date"`
}

var md string
var fm frontmatter

err := mdparse.ParseFile("example.md", &md, &fm)
```