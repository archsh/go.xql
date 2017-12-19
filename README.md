To be done 
=============
To be done.


```golang
package main

import (
    "os"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/archsh/go.xql"
    _ "github.com/archsh/go.xql/dialects/postgres"
)

type School struct {
    xql.Table
    Id      xql.Serial  `xql:"nullable=false,primarykey=true"`
    Name    xql.String  `xql:"length=32,unique=true,index=true"`
    Address xql.String  `xql:"length=256,nullable=true,default=Unknown"`
    Description xql.Text  `xql:"desc,nullable=true"`
}

func (s School) TableName() string {
    return "schools"
}

type Student struct {
    xql.Table
    Id      xql.Serial  `xql:"nullable=false,primarykey=true"`
    Name    xql.String  `xql:"length=32,unique=true,index=true"`
    Number  xql.Integer `xql:"unique=true,index=true"`
    SchoolId xql.Integer    `xql:"foreignkey=@School.Id,nullable=False,index=True"`
    Age     xql.SmallInteger    `xql:"default=0"`
    Telephone xql.String    `xql:"length=32,nullable=true,default=Unknow"`
    Address xql.String  `xql:"length=256,nullable=true,default=Unknown"`
}

func (s Student) TableName() string {
    return "students"
}

func main() {
    db, e := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=cygnuxdb sslmode=disable")
    if nil != e {
        t.Fatal("Connec DB failed:> ", e)
        os.Exit(1)
    }
    session := xql.MakeSession(db, "postgres", true)
    defer session.Close()
    err := session.Create(&School{})
    if nil != err {
        fmt.Errorf("Create table schools failed: %s\n", err)
        os.Exit(1)
    }
    err = session.Create(&Student{})
    if nil != err {
        fmt.Errorf("Create table students failed: %s\n", err)
        os.Exit(1)
    }
    school := &School {
        Name: xql.String("Xinxiu Primary School"),
        Address: xql.String("Xinxiu, Luohu, Shenzhen"),
    }
    student := &Student {
        Name: xql.String("Fangze SHEN"),
        Number: xql.Integer(1),
        SchoolId: school.Id,
        Age: xql.SmallInteger(10),
    }
    if e := session.Insert(school); nil != e {
        fmt.Errorf("Create table students failed: %s\n", err)
        os.Exit(1)
    }
    if e := session.Insert(student); nil != e {
        fmt.Errorf("Create table students failed: %s\n", err)
        os.Exit(1)
    }
    session.Commit()
    session.Close()
}
```
