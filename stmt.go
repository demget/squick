package squick

import (
	"errors"
	"regexp"
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

var reStmtOp = regexp.MustCompile(`\s+\((.*)\)`)

func Parse(s string) (*Stmt, error) {
	tableOps := strings.Split(s, ":")
	if len(tableOps) != 2 {
		return nil, errors.New("make statement should look like table:operations")
	}

	stmt := &Stmt{Table: tableOps[0]}
	ops := strings.Split(tableOps[1], " ")

	for _, op := range ops {
		if !strings.ContainsAny(op, "()") {
			continue
		}

		match := reStmtOp.FindStringSubmatch(op)
		if match == nil {
			return nil, errors.New("make statement operation should look like operation(args,)")
		}

		args := strings.Split(match[1], ",")
		stmt.Operations = append(stmt.Operations, Op{
			Name: op,
			Args: args,
		})
	}

	return stmt, nil
}

func (stmt Stmt) Model() string {
	return strcase.ToCamel(stmt.Table)
}
