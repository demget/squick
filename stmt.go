package squick

import (
	"errors"
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

func Parse(table string, ss []string) (*Stmt, error) {
	if len(ss) == 0 {
		return nil, errors.New("make statement should contain at least one operation")
	}

	stmt := &Stmt{Table: table}
	for _, s := range ss {
		op := Op{Name: s}
		if strings.Contains(s, ":") {
			args := strings.Split(s, ":")[1]
			op.Args = strings.Split(args, ",")
		}

		stmt.Operations = append(stmt.Operations, op)
	}

	return stmt, nil
}

func (stmt Stmt) Model() string {
	return strcase.ToCamel(stmt.Table)
}
