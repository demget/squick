package squick

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
)

const header = "// Code generated by squick at %v."

var (
	//go:embed templates/*
	fs embed.FS

	plur = pluralize.NewClient()
)

type Squick struct {
	tpl *template.Template
}

type Context struct {
	Verbose bool
	Ignore  bool
	DB      *sqlx.DB
	Driver  string
	Package string
	Model   string
	Tags    []string
}

type Column struct {
	DBName   string
	Name     string
	Type     string
	Tags     []string
	Nullable bool
}

func init() {
	strcase.ConfigureAcronym("id", "ID")
}

func New() (*Squick, error) {
	tpl, err := template.
		New("squick").
		Funcs(funcs).
		ParseFS(fs, "templates/*")
	if err != nil {
		return nil, err
	}

	return &Squick{tpl: tpl}, nil
}

func (s *Squick) Init(ctx Context) error {
	if _, err := os.Stat(ctx.Package); os.IsExist(err) {
		return fmt.Errorf("directory %s already exists", ctx.Package)
	}
	if err := os.Mkdir(ctx.Package, 0700); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, "init", ctx); err != nil {
		return err
	}

	data, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/%s.go", ctx.Package, ctx.Package), data, 0700)
}

func (s *Squick) Make(ctx Context, stmt Stmt) error {
	const (
		queryColumns = `
			select column_name, data_type, udt_name
			from information_schema.columns
			where table_name=$1`
		queryPrimary = `
			select kcu.column_name
			from information_schema.table_constraints tco
			join information_schema.key_column_usage kcu
				on kcu.constraint_name = tco.constraint_name
				and kcu.constraint_schema = tco.constraint_schema
				and kcu.constraint_name = tco.constraint_name
			where kcu.table_name=$1 and tco.constraint_type='PRIMARY KEY'`
	)

	var cols []struct {
		Name string `db:"column_name"`
		Type string `db:"data_type"`
		Udt  string `db:"udt_name"`
	}
	if err := ctx.DB.Select(&cols, queryColumns, stmt.Table); err != nil {
		return err
	}

	var primaryKey string
	if err := ctx.DB.Get(&primaryKey, queryPrimary, stmt.Table); err != nil {
		return err // TODO: make primary keys optional
	}

	load := struct {
		Context
		Stmt
		Model        string
		PrimaryKey   string
		Imports      []string
		Dependencies []string
		Columns      []Column
		Blacklist    []string
		ColumnTypes  map[string]string
	}{
		Context:     ctx,
		Stmt:        stmt,
		Model:       plur.Singular(ctx.Model),
		PrimaryKey:  primaryKey,
		Blacklist:   []string{"created_at", "updated_at", primaryKey},
		ColumnTypes: make(map[string]string),
	}

	for _, col := range cols {
		if ctx.Verbose {
			log.Printf("table=%s column=%s type=%s udt=%s\n", stmt.Table, col.Name, col.Type, col.Udt)
		}

		colType, ok := columnTypes[col.Type]
		if !ok {
			if !ctx.Ignore {
				return fmt.Errorf("unsupported %s column type", col.Type)
			} else {
				colType = "interface{}"
			}
		}

		if imp, ok := columnImports[colType]; ok {
			load.Imports = append(load.Imports, imp)
		}
		if dep, ok := columnDependencies[colType]; ok {
			load.Dependencies = append(load.Dependencies, dep)
		}

		colType += udtTypes[col.Udt]
		load.ColumnTypes[col.Name] = colType
		load.Columns = append(load.Columns, Column{
			DBName:   col.Name,
			Name:     toPascalCase(col.Name),
			Type:     colType,
			Tags:     ctx.Tags,
			Nullable: false, // TODO
		})
	}

	for _, op := range stmt.Operations {
		if op.Name == "insert" || op.Name == "update" {
			load.Imports = append(load.Imports, "reflect")
			load.Dependencies = append(load.Dependencies, "github.com/Masterminds/squirrel")
		} else {
			for _, arg := range op.Args {
				var exists bool
				for _, col := range cols {
					if arg == col.Name {
						exists = true
						break
					}
				}
				if !exists {
					log.Fatalf("column %s does not exist", arg)
				}
			}
		}
	}

	var buf bytes.Buffer
	if err := s.tpl.ExecuteTemplate(&buf, "make", load); err != nil {
		return err
	}

	data, err := format.Source(buf.Bytes())
	if err != nil {
		if ctx.Verbose {
			log.Println(buf.String())
		}
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/%s.go", ctx.Package, stmt.Table), data, 0700)
}

var columnTypes = map[string]string{
	"bigint":            "int64",
	"boolean":           "bool",
	"character":         "string",
	"character varying": "string",
	"date":              "time.Time",
	"integer":           "int",
	"real":              "float64",
	"serial":            "int",
	"text":              "string",
	"json":              "types.JSONText",
	"ARRAY":             "pq.",
	"USER-DEFINED":      "string", // TODO: distinguish enums

	"time without time zone":      "time.Time",
	"time with time zone":         "time.Time",
	"timestamp without time zone": "time.Time",
	"timestamp with time zone":    "time.Time",
}

var columnImports = map[string]string{
	"time.Time": "time",
}

var columnDependencies = map[string]string{
	"pq.":            "github.com/lib/pq",
	"types.JSONText": "github.com/jmoiron/sqlx/types",
}

var udtTypes = map[string]string{
	"_varchar": "StringArray",
	"_text":    "StringArray",
	"_int4":    "Int64Array",
	"_int":     "Int32Array",
	"_float4":  "Float64Array",
	"_float":   "Float32Array",
	"_bool":    "BoolArray",
}
