package squick

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type Stmt struct {
	Table      string
	Operations []Op
}

type Op struct {
	Name string
	Args []string
}

func Parse(table string, ss []string) Stmt {
	stmt := Stmt{Table: table}

	for _, s := range ss {
		var op Op

		if strings.Contains(s, ":") {
			opargs := strings.Split(s, ":")
			op.Name = opargs[0]
			op.Args = strings.Split(opargs[1], ",")
		} else {
			op.Name = s
		}

		stmt.Operations = append(stmt.Operations, op)
	}

	return stmt
}

func (stmt Stmt) Model() string {
	return strcase.ToCamel(stmt.Table)
}
