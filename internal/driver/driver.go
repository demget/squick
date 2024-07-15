package driver

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
)

type DB interface {
	// Methods of sqlx.DB struct
	DriverName() string
	MapperFunc(mf func(string) string)
	Rebind(query string) string
	Unsafe() *sqlx.DB
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	MustBegin() *sqlx.Tx
	Beginx() (*sqlx.Tx, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	MustExec(query string, args ...interface{}) sql.Result
	Preparex(query string) (*sqlx.Stmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	MustBeginTx(ctx context.Context, opts *sql.TxOptions) *sqlx.Tx
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Connx(ctx context.Context) (*sqlx.Conn, error)

	// Own methods
	SetTable(string)
	SetDriverImport(string)
	Columns(ignore bool) ([]Column, error)
	PrimaryKey() (string, error)
	Imports() []string
	DriverImport() (_ string, hide bool)
}

type Column struct {
	Name     string
	Type     string
	Nullable bool
}

func New(driver string, uri string) (DB, error) {
	db, err := sqlx.Open(driver, uri)
	if err != nil {
		log.Fatal(err)
	}

	switch driver {
	case "postgres":
		return &postgresDB{
			DB:           db,
			driverImport: "github.com/lib/pq",
		}, nil
	}

	return nil, nil
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
	"numeric":           "string",
	"uuid":              "string",
	"json":              "types.JSONText",

	"time without time zone":      "time.Time",
	"time with time zone":         "time.Time",
	"timestamp without time zone": "time.Time",
	"timestamp with time zone":    "time.Time",
}

var columnImports = map[string]string{
	"time.Time":      "time",
	"types.JSONText": "github.com/jmoiron/sqlx/types",
}
