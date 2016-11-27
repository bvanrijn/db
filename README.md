# db
[![godoc db](http://b.repl.ca/v1/godoc-reference-blue.png)](https://godoc.org/github.com/bvanrijn/db)

A simple database system

## `add.go`
```go
package main

import (
  "fmt"
  
	"github.com/bvanrijn/db"
)

func main() {
  database := db.Database{}
  
  database.Add(db.Record{
    ID:   1,
    URL:  "http://example.com",
    Tags: []string{"example"},
  })
  
  database.Save("data.db")
}
```

`$ go run add.go`

## `serve.go`

```go
package main

import (
  "github.com/bvanrijn/db"
)

func main() {
  database := db.Database{}
  database = database.Load("data.db")
  database.Serve(8000)
}
```

```
$ go run serve.go
$ curl "http://localhost:8000/api?action=search&q=example"
[
  {
    "ID": 1,
    "URL": "http://example.com",
    "Tags": [
      "example"
    ]
  }
]
```
