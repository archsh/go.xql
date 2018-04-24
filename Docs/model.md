# Models

## Basic Table Struct

```go
package xql

type Table struct {
    
}

func (t Table) Construct() 

func (t Table) ConstructSQL() string

func (t Table) Fields()

func (t Table) SetNamespace(string)

func (t Table) Namespace() string

func (t Table) PrimaryKey()

func (t Table) Indexes()

func (t Table) Constraints()

```


## Basic View Struct

```go
package xql

type View struct {
    
}

```


