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

	var set, update bool
	for _, s := range ss {
		op := Op{Name: s}

		if strings.Contains(s, ":") {
			opargs := strings.Split(s, ":")
			op.Name = opargs[0]
			op.Args = strings.Split(opargs[1], ",")
		}

		if !set && op.Name == "set" {
			set = true
		}
		if !update && op.Name == "update" {
			update = true
		}

		stmt.Operations = append(stmt.Operations, op)
	}

	if set && !update {
		stmt.Operations = append(stmt.Operations, Op{Name: "update"})
	}

	return stmt
}

func (stmt Stmt) Model() string {
	return strcase.ToCamel(stmt.Table)
}
