# go.xql Introduction

## Overview


## Levels

- Engine
  Engine is the root resource of DB connection.
- Session
  Session is a temporally instance for short time process.
- QuerySet
  QuerySet is an object contains a specific query.
- Rows
  Rows represents a multi-row result of query.
- Row
  Row represents a single-row result of query.


## Table & Column 

```go
package example

type Person struct {
    Id   int   `xql:"name=id,pk"`
    Name string `xql:"type=varchar,name=name,size=64,unique"`
    Description string `xql:"type=text,nullable"`
}

func (p Person) TableName() string {
    return "persons"
}
```