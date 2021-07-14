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
    squick make select insert update delete
    squick make -table books get:title,author select:year
    squick make -table readers get set:last_visit_at
    squick make -table reader_books insert update`

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
    -table               One specific table to generate model for.

Operations:
    get                  Get by a certain field.      get:author  -> db.BookByAuthor(author)
    select               Select multiple items.       select:year -> db.BooksByYear(year)
    set                  Update a single field.       set:name    -> book.SetTitle(title)
    update               Update entire model.         update      -> db.UpdateBook(id, database.Book{...})
    insert               Insert a model.              insert      -> db.CreateBook(database.Book{...})
    delete               Delete a model.              delete      -> db.DeleteBook(id)
    count                Count aggregations.          count:year  -> db.CountBooksByYear(year)`
