package examples

import (
    "testing"
    "time"
    _ "github.com/lib/pq"
    "github.com/archsh/go.uuid"
    "github.com/archsh/go.xql"
    "github.com/archsh/go.xql/dialects/postgres"
    "os"
    "fmt"
)

type Category struct {
    Id          string     `json:"id" xql:"type=uuid,pk,default=uuid_generate_v4()"`
    Name        string     `json:"name" xql:"size=24,unique,index"`
    Tags        postgres.StringArray `json:"tags" xql:"size=32,nullable"`
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
    Attributes  postgres.HSTORE `json:"attributes" xql:"nullable"`
    Scores      postgres.RealArray `json:"scores" xql:"nullable"`
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
var session *xql.Session

func TestMain(m *testing.M) {
    var retCode int
    var e error
    session, e = prepare()
    if nil != e {
        fmt.Println("> Prepare failed:>", e)
        os.Exit(-1)
    }
    e = destroy_tables(session, MovieCrew, MovieCategory)
    if nil != e {
        fmt.Println("> Destroy Tables failed:>", e)
        os.Exit(-1)
    }
    e = create_tables(session, MovieCategory, MovieCrew)
    if nil != e {
        fmt.Println("> Create Tables failed:>", e)
        os.Exit(-1)
    }
    if nil == e {
        retCode = m.Run()
    }else{
        retCode = -1
    }
    os.Exit(retCode)
}


func prepare() (*xql.Session, error) {
    engine, e := xql.CreateEngine("postgres",
        "host=localhost port=5432 user=postgres password=postgres dbname=test sslmode=disable")
    if nil != e {
        return nil, e
    }
    sess := engine.MakeSession()
    return sess, nil
}

func create_tables(s *xql.Session, tables... *xql.Table) error {
    for _, t := range tables {
        e := s.Create(t)
        if nil != e {
            return e
        }
    }

    return nil
}

func destroy_tables(s *xql.Session, tables... *xql.Table) error {
    for _, t := range tables {
        e := s.Drop(t, true)
        if nil != e {
            return e
        }
    }
    return nil
}

func TestQuerySet_Insert(t *testing.T) {
    t1 := time.Now()
    c := Category{Id:uuid.NewV4().String(), Name:"Test"}
    c.Tags = []string{"Star", "Actor"}
    c1 := Crew{Id: uuid.NewV4().String(), FullName: "Tom Cruse", Region: "US", Age: 19, CategoryId:c.Id}
    c1.Attributes = make(map[string]interface{})
    c1.Attributes["skill"] = "Good"
    c1.Attributes["score"] = "99.5"
    c1.Scores = []float32{99,96,93}
    c2 := Crew{Id: uuid.NewV4().String(), FullName: "Hue Jackman", Region: "US", Age:21, CategoryId:c.Id}
    c2.Attributes = make(map[string]interface{})
    c2.Attributes["skill"] = "Normal"
    c2.Attributes["score"] = "79.5"
    c2.Scores = []float32{89,76,83}
    if n, e := session.Query(MovieCategory).Insert(&c); nil != e {
        t.Fatal("Insert failed:> ", e)
    }else{
        t.Log("Inserted Category: ", n, c)
    }
    if n, e := session.Query(MovieCrew).Insert(&c1, &c2); nil != e {
        t.Fatal("Insert failed:> ", e)
    }else{
        t.Log("Inserted Crew: ", n, c1,c2)
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_One(t *testing.T) {
    t1 := time.Now()
    crew := Crew{}
    if e := session.Query(MovieCrew).One().Scan(&crew); nil != e {
        t.Fatal("Qery One failed:>", e)
    }else{
        t.Log("Queried One:>", crew)
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Scan(t *testing.T) {
    t1 := time.Now()
    if r, e :=session.Query(MovieCrew).Filter(map[string]interface{}{"region":"US"}).All(); nil != e {
        t.Fatal("Query all failed:>", e)
    }else{
        defer r.Close()
        for r.Next() {
            crew := Crew{}
            if e = r.Scan(&crew); nil != e {
                t.Fatal("Scan failed:>", e)
            }else{
                t.Log("Scanned Crew:>", crew)
            }
        }
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Update(t *testing.T) {
    t1 := time.Now()
    if n, e := session.Query(MovieCrew).Update(map[string]interface{}{"age": 30}); nil != e {
        t.Fatal("Update failed:>", e)
    }else{
        t.Log("Updated rows:>", n)
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Delete(t *testing.T) {
    t1 := time.Now()
    if n, e := session.Query(MovieCrew).Filter(map[string]interface{}{"full_name":"Tom Cruse"}).Delete(); nil != e {
        t.Fatal("Delete failed:>", e)
    }else{
        t.Log("Deleted rows:>", n)
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}


