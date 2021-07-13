package main

const help = `
Usage: squick <command> [arguments]

Squick generates highly idiomatic Go code to interact with SQL database.

Commands:
    init                Initializes a database package.
    make                Generates code for the specified table.

Drivers:
    postgres
    mysql                (unadapted for now)
    sqlite3              (unadapted for now)

Example:
    export SQUICK_DRIVER="postgres"
    export SQUICK_URL="host=localhost sslmode=disable user=... password=... dbname=..."
    squick init database
    squick make base
    squick make -table books get:id,title,author count
    squick make -table reader_books base:select,insert`

const helpInit = `
Usage: squick init [options] <package>

Squick init initializes a database package.

Options:
    -force               Forces the recreation of the package (wipes out the entire directory).`

const helpMake = `
Usage: squick make [options] <table:operations>

Squick make generates code for the specified table with specified operations.

Options:
    -name                Generated model name, turned into singular PascalCase by default.
    -tags                Additional tags to define for the fields, json only by default.

Operations:
    base(ops)            Primary operations: select (by a primary key), insert, update, delete.
    get(fields)          Get operations. To use certain fields list them inside the parens.
    set(fields)          Set operations. To use certain fields list them inside the parens.
    count(fields)        Count aggregation. To group by certain fields list them inside the parens.`
