# Squick
> `go install github.com/demget/squick/cmd/squick`

_Note: The code is a bit dirty at the moment, but someday I will spruce it up..._

- [Bootstrap](#bootstrap)
- [Models](#models)
- [Roadmap](#roadmap)

## Bootstrap

```
$ squick init database
```

This will produce a package `database`, which holds initial `database.go` file and used to store further generated models and query functions.

```go
db, err := database.Open(os.Getenv("DB_URL"))
if err != nil {
	panic(err)
}

// Now we can use wrapped db to interact with database
// db.UserByID(...)
// db.InsertUser(database.User{...})
```

## Models

Let's imagine we're developing a backend for library application. In our Postgres database we have simple `books`, `readers` and `reader_books` tables, and the goal is to quickly get some fancy query wrappers to deal with basic CRUD operations. We want to select an information about single or multiple books as well as change some reader's data and track taken books by him.

```sql
create table books (
    created_at      timestamp,
    id              serial primary key,
    isbn            varchar(17) not null,
    title           varchar(256) not null,
    author          varchar(256) not null,
    year            int,
    pages           int
);

create table readers (
    created_at      timestamp,
    id              serial primary key,
    full_name       varchar(256) not null,
    age             int,
    last_visit_at   timestamp
);

create table reader_books (
    took_at         timestamp,
    id              serial primary key,
    reader_id       int,
    book_id         int,
    returned        bool
);
```

`get` operation creates queries that retrieve information about single book by the listed fields. We also need to filter by a year, fetching multiple books.

```
$ squick make -table books get:title,author select:year
```

```go
book, _ := 
	db.BookByID(1)
	db.BookByTitle("Clockwork Orange") 
	db.BookByAuthor("Jack London")
```

```go
books, _ := db.BooksByYear(1995)
for _, book := range books {
	// every single book published in 1995
}
```

Now, we need an update function to actualize last visit time of the person came to the library. Setters are only available in the scope of the model.

```
$ squick make -table readers get set:last_visit_at
```

```go
reader, _ := db.ReaderByID(1)
reader.SetLastVisitAt(time.Now())
```

And, finally, an insert operation to track taken books. In this example we use `update` operation instead of simple `set:returned`, just to show the difference (`update` operation allows to update several fields simultaneously).

```
$ squick make -table reader_books insert update
```

```go
id, _ := db.InsertReaderBook(database.Book{
	ReaderID: 1,
	BookID:   20,
})

db.UpdateReaderBook(id, database.Book{
	Returned: true,
})
```

## Roadmap

Here are some plans on this project. It's fair to say that `squick` is __unfinished__. It could have a lot more features and be more flexible, but for now it satisfies my requests and significantly helps me with the CRUD routine in the way I want it.

- Clean up the code and templates
- Make current basic operations more flexible and customizable
- Add support for other drivers besides postgres
- Cover with tests, add them to generated packages as well
- Add mock generation for testing purposes

And some small missed features:
- Deal with nullable fields
- Empty `-table` option should mean to process every table in the database