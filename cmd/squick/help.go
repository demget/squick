package main

const help = `
Usage: squick <command> [arguments]

Squick generates highly idiomatic Go code to interact with SQL databases.

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
    squick make -table books get:title,author select:year
    squick make -table readers get set:last_visit_at
    squick make -table reader_books insert update delete`

const helpInit = `
Usage: squick init [options] <package>

Squick init initializes a database package.

Options:
    -force               Forces the recreation of the package (wipes out the entire directory).
    -max-open            Sets the maximum number of connections in the idle connection pool.
    -max-idle            Sets the maximum number of open connections to the database.`

const helpMake = `
Usage: squick make [options] <table:operations>

Squick make generates code for the specified table with specified operations.

Options:
    -v                   Enables verbose debug output.
    -ignore              Ignores unsupported column, interface{} type will be used instead.
    -nopk                Ignores primary key absence.
    -table               One specific table to generate model for, required.
    -name                Custom model name, singular of table name by default.
    -tags                Additional tags to define for the fields, json only by default.
    -updated             Indicates a field, which has to be set to the latest time on each updating query.

Operations:
    get                  Get by a certain field.      get:author  -> db.BookByAuthor(author)
    select               Select multiple items.       select:year -> db.BooksByYear(year)
    set                  Update a single field.       set:title   -> book.SetTitle(title)
    update               Update an entire model.      update      -> db.UpdateBook(id, database.Book{...})
    insert               Insert a model.              insert      -> db.InsertBook(database.Book{...})
    delete               Delete a model.              delete      -> db.DeleteBook(id)
    count                Count by a field.            count:year  -> db.CountBooksByYear(year)`
