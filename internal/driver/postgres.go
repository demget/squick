package driver

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type postgresDB struct {
	*sqlx.DB

	table        string
	driverImport string
	driverUsed   bool
	imports      []string
}

func (p *postgresDB) SetTable(name string) {

	p.table = name
}

func (p *postgresDB) PrimaryKey() (primaryKey string, _ error) {
	const queryPrimary = `
			select kcu.column_name
			from information_schema.table_constraints tco
			join information_schema.key_column_usage kcu
				on kcu.constraint_name = tco.constraint_name
				and kcu.constraint_schema = tco.constraint_schema
				and kcu.constraint_name = tco.constraint_name
			where kcu.table_name=$1 and tco.constraint_type='PRIMARY KEY'`
	if err := p.Get(&primaryKey, queryPrimary, p.table); err != nil {
		return "", err
	}

	return primaryKey, nil
}

func (p *postgresDB) Columns(ignore bool) ([]Column, error) {
	const queryColumns = `
			select column_name, data_type, udt_name, is_nullable
			from information_schema.columns
			where table_name=$1`

	var cols []struct {
		Name     string `db:"column_name"`
		Type     string `db:"data_type"`
		Udt      string `db:"udt_name"`
		Nullable string `db:"is_nullable"`
	}
	if err := p.Select(&cols, queryColumns, p.table); err != nil {
		return nil, err
	}

	result := make([]Column, len(cols))
	imports := make(map[string]struct{})

	for i, col := range cols {
		result[i].Name = col.Name

		switch col.Type {
		// postgres types
		case "ARRAY":
			p.driverUsed = true
			result[i].Type = "pg." + udtTypes[col.Udt]
		case "USER-DEFINED":
			result[i].Type = "string"
		default:
			if im, ok := columnImports[columnTypes[col.Type]]; ok {
				imports[im] = struct{}{}
			}
			colType, ok := columnTypes[col.Type]
			if !ok && !ignore {
				return nil, fmt.Errorf("supported column type: %s", col.Type)
			}

			result[i].Type = colType
		}

		if col.Type != "json" && col.Type != "ARRAY" {
			result[i].Nullable = col.Nullable == "YES"
		}
	}

	for imp := range imports {
		p.imports = append(p.imports, imp)
	}

	return result, nil
}

func (p *postgresDB) SetDriverImport(i string) {
	p.driverImport = i
}

func (p *postgresDB) Imports() []string {
	return p.imports
}

func (p *postgresDB) DriverImport() (string, bool) {
	return p.driverImport, !p.driverUsed
}

var udtTypes = map[string]string{
	"_varchar": "StringArray",
	"_text":    "StringArray",
	"_int2":    "Int32Array",
	"_int4":    "Int32Array",
	"_int8":    "Int64Array",
	"_int":     "Int32Array",
	"_float4":  "Float64Array",
	"_float":   "Float32Array",
	"_bool":    "BoolArray",
}
