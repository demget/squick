package squick

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
	"golang.org/x/mod/modfile"
)

var funcs = template.FuncMap{
	"header":  func() string { return fmt.Sprintf(header, time.Now().Format(time.RFC3339)) },
	"command": func() string { return "// squick " + strings.Join(os.Args[1:], " ") },
	"isint":   func(t string) bool { return strings.HasPrefix(t, "int") },

	"package": packageFromPath,
	"modpath": modulePath,
	"in":      inStrings,
	"camel":   toCamelCase,
	"pascal":  strcase.ToCamel,
	"plural":  plur.Plural,
}

func packageFromPath(pkg string) string {
	if !strings.Contains(pkg, "/") {
		return pkg
	}
	return pkg[strings.LastIndex(pkg, "/")+1:]
}

func modulePath() string {
	data, _ := ioutil.ReadFile("go.mod")
	return modfile.ModulePath(data)
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
