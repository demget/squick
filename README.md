# Squick

## Initialize database package

`$ squick init database`

This will produce a package `database`, which holds initial `database.go` file and used to store further generated models and query functions.

```go
db, err := database.Open(os.Getenv("DB_URL"))
if err != nil {
	panic(err)
}

// Now we can use wrapped db to interact with database
// db.User(...)
// db.CreateUser(database.User{...})
```

## Generate model and functions

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
    reader_id       int unique,
    book_id         int unique,
    returned        bool
);
```

We specify `get` operation to have a look on functions help to retrieve information about single book by a particular filter rule. Then we need to filter by a year, but fetching multiple books.

`$ squick make -table books get:title,author select:year`

```go
book, _ := 
	db.Book(1) // by id (primary key)
	db.BookByTitle("Clockwork Orange") 
	db.BookByAuthor("Jack London")
```

```go
books, _ := db.BooksByYear(1995)
for _, book := range books {
	// every single book published in 1995
}
```

Now, we need an update function to actualize last visit time of the person came to the library. Setters are only available in the scope of model.

`$ squick make -table readers get set:last_visit_at`

```go
reader, _ := db.Reader(readerID)
reader.SetLastVisitAt(time.Now())
```

And, finally, an insert operation to track taken books. In this example we use `update` operation instead of simple `set:returned`, just to show the difference.

`$ squick make -table reader_books insert update`

```go
id, _ := db.CreateReaderBook(database.Book{
	ReaderID: 1,
	BookID:   20,
})

db.UpdateReaderBook(id, database.Book{
	Returned: true,
})
```
