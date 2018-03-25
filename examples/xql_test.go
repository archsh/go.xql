package examples

import (
    "testing"
    "time"
    _ "github.com/lib/pq"
    "github.com/archsh/go.xql"
    "github.com/archsh/go.xql/dialects/postgres"
    "os"
    "fmt"
    "database/sql/driver"
)

type School struct {
    Id          int                  `json:"id" xql:"type=serial,pk"`
    Name        string               `json:"name" xql:"size=24,unique,index"`
    Tags        postgres.StringArray `json:"tags" xql:"size=32,nullable"`
    Description string               `json:"description"  xql:"name=desc,type=text,size=24,nullable=false,default=''"`
}

func (c School) TableName() string {
    return "schools"
}

type Character struct {
    postgres.JSON
    Attitude string
    Height   int
    Weight   int
}

//func (j Character) Declare(props xql.PropertySet) string {
//    return "JSONB"
//}

func (j *Character) Scan(value interface{}) error {
    return postgres.JSONB_Scan(j, value)
}

func (j Character) Value() (driver.Value, error) {
    return postgres.JSONB_Value(j)
}

type People struct {
    Id          int        `json:"id" xql:"type=serial,pk"`
    FullName    string     `json:"fullName" xql:"size=80,unique=true,index=true"`
    FirstName   string     `json:"firstName" xql:"size=24,default=''"`
    MiddleName  string     `json:"middleName" xql:"size=24,default=''"`
    LastName    string     `json:"lastName" xql:"size=24,default=''"`
    Region      string     `json:"region"  xql:"size=24,nullable=true"`
    Age         int        `json:"age" xql:"check=(age>18)"`
    SchoolId    int        `json:"schoolId"  xql:"type=integer,fk=schools.id,ondelete=CASCADE"`
    Description string     `json:"description"  xql:"name=desc,type=text,size=24,default=''"`
    Created     *time.Time `json:"created"  xql:"type=timestamp,default=Now()"`
    Updated     *time.Time `json:"Updated"  xql:"type=timestamp,default=Now()"`
}

type Teacher struct {
    People
    Degree string `json:"degree" xql:"size=64,default=''"`
}

func (t Teacher) TableName() string {
    return "teachers"
}

type Student struct {
    People
    Grade      string             `json:"grade" xql:"size=32,default=''"`
    Attributes postgres.HSTORE    `json:"attributes" xql:"nullable"`
    Scores     postgres.RealArray `json:"scores" xql:"nullable"`
    Character  Character          `json:"character" xql:"nullable"`
}

func (c Student) TableName() string {
    return "students"
}

var StudentTable = xql.DeclareTable(&Student{})
var TeacherTable = xql.DeclareTable(&Teacher{})
var SchoolTable = xql.DeclareTable(&School{})
var session *xql.Session

func TestMain(m *testing.M) {
    var retCode int
    var e error
    session, e = prepare()
    if nil != e {
        fmt.Println("> Prepare failed:>", e)
        os.Exit(-1)
    }
    //e = destroy_tables(session, StudentTable, TeacherTable, SchoolTable)
    //if nil != e {
    //    fmt.Println("> Destroy Tables failed:>", e)
    //    os.Exit(-1)
    //}
    e = create_tables(session, SchoolTable, TeacherTable, StudentTable)
    if nil != e {
        fmt.Println("> Create Tables failed:>", e)
        os.Exit(-1)
    }
    if nil == e {
        retCode = m.Run()
    } else {
        retCode = -1
    }
    os.Exit(retCode)
}

func prepare() (*xql.Session, error) {
    engine, e := xql.CreateEngine("postgres",
        "host=postgresql port=5432 user=postgres password=postgres dbname=test sslmode=disable")
    if nil != e {
        return nil, e
    }
    if e := postgres.Initialize_HSTORE(engine.DB(), "public"); nil != e {
        return nil, e
    }
    sess := engine.MakeSession()
    return sess, nil
}

func create_tables(s *xql.Session, tables ... *xql.Table) error {
    for _, t := range tables {
        e := s.Create(t)
        if nil != e {
            return e
        }
    }

    return nil
}

func destroy_tables(s *xql.Session, tables ... *xql.Table) error {
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
    c := School{Name: "Xinxiu Primary School", Description: "Xinxiu"}
    c.Tags = []string{"Primary", "Luohu", "Shenzhen", "Public"}
    c1 := Student{People: People{FullName: "Tom Cruse", Region: "US", Age: 19, SchoolId: c.Id}}
    c1.Attributes = make(map[string]interface{})
    c1.Attributes["skill"] = "Good"
    c1.Attributes["score"] = "99.5"
    c1.Scores = []float32{99, 96, 93}
    c1.Character.Attitude = "OK"
    c1.Character.Height = 172
    c1.Character.Weight = 68
    c2 := Student{People: People{FullName: "Hue Jackman", Region: "US", Age: 21, SchoolId: c.Id}}
    c2.Attributes = make(map[string]interface{})
    c2.Attributes["skill"] = "Normal"
    c2.Attributes["score"] = "79.5"
    c2.Scores = []float32{89, 76, 83}
    c2.Character.Attitude = "Good"
    c2.Character.Height = 192
    c2.Character.Weight = 88
    var id int
    if e := session.Query(SchoolTable).InsertWithInsertedId(&c, "id", &id); nil != e {
        t.Fatal("Insert failed:> ", e)
    } else {
        t.Log("Inserted School: ", id, c)
    }
    c1.SchoolId = id
    c2.SchoolId = id
    if n, e := session.Query(StudentTable).Insert(&c1, &c2); nil != e {
        t.Fatal("Insert failed:> ", e)
    } else {
        t.Log("Inserted Student: ", n, c1, c2)
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_One(t *testing.T) {
    t1 := time.Now()
    crew := Student{}
    if e := session.Query(StudentTable).One().Scan(&crew); nil != e {
        t.Fatal("Qery One failed:>", e)
    } else {
        t.Log("Queried One:>", crew)
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Scan(t *testing.T) {
    t1 := time.Now()
    if r, e := session.Query(StudentTable).Filter(map[string]interface{}{"region": "US"}).All(); nil != e {
        t.Fatal("Query all failed:>", e)
    } else {
        defer r.Close()
        for r.Next() {
            crew := Student{}
            if e = r.Scan(&crew); nil != e {
                t.Fatal("Scan failed:>", e)
            } else {
                t.Log("Scanned Student:>", crew)
            }
        }
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Update(t *testing.T) {
    t1 := time.Now()
    if n, e := session.Query(StudentTable).Update(map[string]interface{}{"age": 30}); nil != e {
        t.Fatal("Update failed:>", e)
    } else {
        t.Log("Updated rows:>", n)
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Delete(t *testing.T) {
    t1 := time.Now()
    if n, e := session.Query(StudentTable).Filter(map[string]interface{}{"full_name": "Tom Cruse"}).Delete(); nil != e {
        t.Fatal("Delete failed:>", e)
    } else {
        t.Log("Deleted rows:>", n)
    }
    t.Log("Time spent:> ", time.Now().Sub(t1))
}
