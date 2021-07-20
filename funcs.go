package squick

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

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

func toCamelCase(s string) string {
	if s == "id" {
		return s
	}
	s = strcase.ToLowerCamel(s)
	s = strings.ReplaceAll(s, "Id", "ID")
	return s
}

func toPascalCase(s string) string {
	s = strcase.ToCamel(s)
	s = strings.ReplaceAll(s, "Id", "ID")
	return s
}
