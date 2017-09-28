To SQLAlchemy 
=============
To be done.


```golang
package main

import (
    "fmt"
    _ "github.com/lib/pq"
    "github.com/archsh/go.xql"
    _ "github.com/archsh/go.xql/dialects/postgres"
)
type School struct {
    Id int `xql:"Column('id', Integer, Sequence('school_id_seq'), nullable=False, primary_key=True, default=1)"`
    Name string `xql:"Column('name', Unicode(64), nullable=False, unique=True, default='')"`
    Address string `xql:"Column('address', Unicode(256), nullable=True)"`
    Description string `xql:"Column('desc', Text, nullable=True)"`
}
type Student struct {
    Id int `xql:"Column('id', Integer, Sequence('student_id_seq'), nullable=False, primary_key=True, default=1)"`
    Name string `xql:"Column('name', Unicode(64), nullable=False, unique=True, default='')"`
    Number int `xql:"Column('number', Integer, nullable=False, unique=True, default=1)"`
    SchoolId int `xql:"Column('school_id', Integer, ForeignKey('schools.id', ondelete='CASCADE', onupdate='CASCADE'), nullable=False)"`
    Age int `xql:"Column('age', SmallInteger, nullable=False, default=1)"`
    Telephone string `xql:"Column('telephone', String(32), nullable=True)"`
    Address string `xql:"Column('address', Unicode(256), nullable=True)"`
}

var SchoolTable = xql.DeclareTable("schools", &School{})
var StudentTable = xql.DeclareTable("students", &Student{})

func main() {
    db, e := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=cygnuxdb sslmode=disable")
    if nil != e {
        t.Fatal("Connec DB failed:> ", e)
    }
    session := xql.MakeSession(db, "postgres", true)
    defer session.Close()
    err := session.Create(SchoolTable)
    if nil != err {
        fmt.Errorf("Create table schools failed: %s\n", err)
    }
    err = session.Create(StudentTable)
    if nil != err {
        fmt.Errorf("Create table students failed: %s\n", err)
    }
}
```
