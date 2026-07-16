package mdparse

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseFile(file string, markdown *string, v any) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return Parse(string(content), markdown, v)
}

func Parse(input string, markdown *string, v any) error {
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

func extract(data []byte) (hasFrontmatter bool, frontmatter string, markdown string, err error) {
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM

	sc := bufio.NewScanner(strings.NewReader(string(data)))
	// falls du große Dateien/Frontmatter hast:
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	// Datei muss mit --- (als Zeile) beginnen
	if !sc.Scan() {
		return false, "", "", nil
	}
	if strings.TrimRight(sc.Text(), "\r") != "---" {
		return false, "", "", nil
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
	return false, "", "", nil
}
