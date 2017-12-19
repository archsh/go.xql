package examples

import (
    "testing"
    "time"
    _ "github.com/lib/pq"
    "github.com/archsh/go.uuid"
    "github.com/archsh/go.xql"
    _ "github.com/archsh/go.xql/dialects/postgres"
)

type Category struct {
    Id          string     `json:"id" xql:"type=uuid,pk,default=uuid_generate_v4()"`
    Name        string     `json:"name" xql:"size=24,unique,index"`
    Description string     `json:"description"  xql:"name=desc,type=text,size=24,nullable=false"`
}

func (c Category) TableName() string {
    return "categories"
}

type Crew struct {
    Id          string     `json:"id" xql:"type=uuid,primarykey,default=uuid_generate_v4()"`
    FullName    string     `json:"fullName" xql:"size=80,unique=true,index=true"`
    FirstName   string     `json:"firstName" xql:"size=24"`
    MiddleName  string     `json:"middleName" xql:"size=24"`
    LastName    string     `json:"lastName" xql:"size=24"`
    Region      string     `json:"region"  xql:"size=24,nullable=true"`
    Age         int         `json:"age" xql:"check=(age>18)"`
    CategoryId  string     `json:"categoryId"  xql:"type=uuid,fk=categories.id,ondelete=CASCADE"`
    Description string     `json:"description"  xql:"name=desc,type=text,size=24"`
    Created     *time.Time `json:"created"  xql:"type=timestamp,default=Now()"`
    Updated     *time.Time `json:"Updated"  xql:"type=timestamp,default=Now()"`
}

func (c Crew) TableName() string {
    return "crews"
}

var MovieCrew = xql.DeclareTable(&Crew{}, )
var MovieCategory = xql.DeclareTable(&Category{})

func TestCreateEngine(t *testing.T) {
    t1 := time.Now()
    engine, e := xql.CreateEngine("postgres",
        "host=localhost port=5432 user=postgres password=postgres dbname=test sslmode=disable")
    if nil != e {
        t.Fatal("Connec DB failed:> ", e)
    }
    t.Log("MovieCrew:> ", MovieCrew)
    sess := engine.MakeSession()
    sess.Create(MovieCrew)
    t.Log("Time spent:> ", time.Now().Sub(t1))

}

func TestQuerySet_Insert(t *testing.T) {
    t1 := time.Now()
    engine, e := xql.CreateEngine("postgres",
        "host=localhost port=5432 user=postgres password=postgres dbname=test sslmode=disable")
    if nil != e {
        t.Fatal("Connec DB failed:> ", e)
    }
    t.Log("MovieCrew:> ", MovieCrew)
    sess := engine.MakeSession()
    e = sess.Drop(MovieCrew, true)
    if nil != e {
        t.Fatal("Failed to drop table:>", e)
        return
    }
    e = sess.Drop(MovieCategory, true)
    if nil != e {
        t.Fatal("Failed to drop table:>", e)
        return
    }
    e = sess.Create(MovieCategory)
    if nil != e {
        t.Fatal("Failed to create table:>", e)
        return
    }
    e = sess.Create(MovieCrew)
    if nil != e {
        t.Fatal("Failed to create table:>", e)
        return
    }
    c := Category{Id:uuid.NewV4().String(), Name:"Test"}
    c1 := Crew{Id: uuid.NewV4().String(), FullName: "Tom Cruse", Region: "US", Age: 19, CategoryId:c.Id}
    c2 := Crew{Id: uuid.NewV4().String(), FullName: "Hue Jackman", Region: "US", Age:16, CategoryId:c.Id}
    if n, e := sess.Query(MovieCategory).Insert(&c); nil != e {
        t.Fatal("Insert failed:> ", e)
    }else{
        t.Log("Inserted Category: ", n, c)
    }
    if n, e := sess.Query(MovieCrew).Insert(&c1, &c2); nil != e {
        t.Fatal("Insert failed:> ", e)
    }else{
        t.Log("Inserted Crew: ", n, c1,c2)
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Update(t *testing.T) {
    t1 := time.Now()

    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Delete(t *testing.T) {
    t1 := time.Now()

    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_One(t *testing.T) {
    t1 := time.Now()

    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Scan(t *testing.T) {
    t1 := time.Now()

    t.Log("Time spent:> ", time.Now().Sub(t1))
}
