go.xql, another ORM for golang 
==============================
To be done.


See [example](examples/xql_test.go)
```go
package main
import "github.com/archsh/go.xql"
var table = xql.Table(
	xql.Column("id",),
	xql.Index("id"),
	xql.Constraint(""),
	xql.PrimaryKey("id"),
	xql.Comment(""),
	)
```