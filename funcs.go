package squick

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
	"unicode"

	"github.com/iancoleman/strcase"
)

var funcs = template.FuncMap{
	"header":  func() string { return fmt.Sprintf(header, time.Now().Format(time.RFC3339)) },
	"command": func() string { return "// squick " + strings.Join(os.Args[1:], " ") },
	"isint":   func(t string) bool { return strings.HasPrefix(t, "int") },

	"package": packageFromPath,
	"in":      inStrings,
	"camel":   toCamelCase,
	"pascal":  toPascalCase,
	"plural":  plur.Plural,
}

func init() {
	// TODO: add an ability to specify acronyms (v1)
	strcase.ConfigureAcronym("id", "ID")
}

func packageFromPath(pkg string) string {
	if !strings.Contains(pkg, "/") {
		return pkg
	}
	return pkg[strings.LastIndex(pkg, "/")+1:]
}

func inStrings(a []string, v string) bool {
	for _, s := range a {
		if s == v {
			return true
		}
	}
	return false
}

func replaceAcronyms(s, a string) string {
	var (
		word  string
		words []string
	)
	for _, r := range []rune(s) {
		if unicode.IsUpper(r) {
			words = append(words, word)
			word = ""
		}
		word += string(r)
	}

	words = append(words, word)
	for i, word := range words {
		if word == a {
			words[i] = strings.ToUpper(word)
		}
	}

	return strings.Join(words, "")
}

func toCamelCase(s string) string {
	if s == "id" {
		return s
	}
	s = strcase.ToLowerCamel(s)
	s = replaceAcronyms(s, "Id")
	return s
}

func toPascalCase(s string) string {
	s = strcase.ToCamel(s)
	s = replaceAcronyms(s, "Id")
	return s
}
