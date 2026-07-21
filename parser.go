package mdparse

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadFile is like Read, but it reads the input from a file instead of a string.
func ReadFile(file string, markdown *string, v any) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return Read(string(content), markdown, v)
}

// Read reads the frontmatter and markdown content from the given input string.
// It returns the markdown content unmodified and unmarshals the frontmatter into
// the provided struct. The input string is expected to have the frontmatter enclosed
// between two lines containing only '---'. If no frontmatter is found, the markdown
// content will be returned as-is, and the provided struct will remain unchanged.
func Read(input string, markdown *string, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	// remove any UTF-8 BOM if present
	hasFm, fm, md, err := extract([]byte(input))
	if err != nil {
		return err
	}

	*markdown = md

	// lastly, unmarshal the frontmatter into the provided struct
	if hasFm {
		err = yaml.Unmarshal([]byte(fm), v)
	}

	return err
}

func Parse(input string) ([]Paragraph, error) {
	panic("implement me")
}

func extract(data []byte) (hasFrontmatter bool, frontmatter string, markdown string, err error) {
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM
	plainMarkdown := strings.TrimRight(string(data), "\r\n")

	sc := bufio.NewScanner(strings.NewReader(string(data)))
	// falls du große Dateien/Frontmatter hast:
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	// Datei muss mit --- (als Zeile) beginnen
	if !sc.Scan() {
		return false, "", "", nil
	}
	if strings.TrimRight(sc.Text(), "\r") != "---" {
		return false, "", plainMarkdown, nil
	}

	// Frontmatter bis zum nächsten --- sammeln
	var fm strings.Builder
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if line == "---" {
			// Der Rest ist der eigentliche Markdown-Inhalt
			var md strings.Builder
			first := true
			for sc.Scan() {
				if !first {
					md.WriteByte('\n')
				}
				first = false
				md.WriteString(strings.TrimRight(sc.Text(), "\r"))
			}
			return true, fm.String(), md.String(), nil
		}
		fm.WriteString(line)
		fm.WriteByte('\n')
	}

	// Kein schließendes ---
	return false, "", plainMarkdown, nil
}
