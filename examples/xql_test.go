package examples

import (
	"database/sql/driver"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/archsh/go.xql"
	"github.com/archsh/go.xql/dialects/postgres"
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
	return postgres.JsonbScan(j, value)
}

func (j Character) Value() (driver.Value, error) {
	return postgres.JsonbValue(j)
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
	Grade      string               `json:"grade" xql:"size=32,default=''"`
	Attributes postgres.HSTORE      `json:"attributes" xql:"nullable"`
	Scores     postgres.DoubleArray `json:"scores" xql:"nullable"`
	Character  Character            `json:"character" xql:"nullable"`
}

func (c Student) TableName() string {
	return "students"
}

var StudentTable = xql.DeclareTable(&Student{})
var TeacherTable = xql.DeclareTable(&Teacher{})
var SchoolTable = xql.DeclareTable(&School{})
var session *xql.Session
var schoolId int

func TestMain(m *testing.M) {
	var retCode int
	var e error

	if session, e = prepare(); nil != e {
		fmt.Println("> Prepare failed:>", e)
		os.Exit(-1)
	}

	e = createTables(session, SchoolTable, TeacherTable, StudentTable)
	if nil != e {
		fmt.Println("> Create Tables failed:>", e)
		os.Exit(-1)
	}

	c := School{Name: "Xinxiu Primary School", Description: "Xinxiu"}
	c.Tags = []string{"Primary", "Luohu", "Shenzhen", "Public"}
	if e := session.Table(SchoolTable).InsertWithInsertedId(&c, "id", &schoolId); nil != e {
		fmt.Println("Insert failed:> ", e)
		os.Exit(-1)
	} else {
		fmt.Println("Inserted School: ", schoolId, c)
	}
	retCode = m.Run()

	if e := destroyTables(session, StudentTable, TeacherTable, SchoolTable); nil != e {
		fmt.Println("> Destroy Tables failed:>", e)
		os.Exit(-1)
	}
	os.Exit(retCode)
}

func prepare() (*xql.Session, error) {
	engine, e := xql.CreateEngine("postgres",
		"host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable")
	if nil != e {
		return nil, e
	}
	if e := postgres.InitializeHSTORE(engine.DB(), "public"); nil != e {
		return nil, e
	}
	sess := engine.MakeSession()
	return sess, nil
}

func createTables(s *xql.Session, tables ...*xql.Table) error {
	for _, t := range tables {
		e := s.Create(t)
		if nil != e {
			return e
		}
	}

	return nil
}

func destroyTables(s *xql.Session, tables ...*xql.Table) error {
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

	c1 := Student{People: People{FullName: "Tom Cruse", Region: "US", Age: 19, SchoolId: schoolId}}
	c1.Attributes = make(map[string]interface{})
	c1.Attributes["skill"] = "Good"
	c1.Attributes["score"] = "99.5"
	c1.Scores = []float64{99, 96, 93}
	c1.Character.Attitude = "OK"
	c1.Character.Height = 172
	c1.Character.Weight = 68
	c2 := Student{People: People{FullName: "Hue Jackman", Region: "US", Age: 21, SchoolId: schoolId}}
	c2.Attributes = make(map[string]interface{})
	c2.Attributes["skill"] = "Normal"
	c2.Attributes["score"] = "79.5"
	c2.Scores = []float64{89, 76, 83}
	c2.Character.Attitude = "Good"
	c2.Character.Height = 192
	c2.Character.Weight = 88
	c1.SchoolId = schoolId
	c2.SchoolId = schoolId
	if n, e := session.Table(StudentTable).Insert(&c1, &c2); nil != e {
		t.Fatal("Insert failed:> ", e)
	} else {
		t.Log("Inserted Students: ", n)
		for _, c := range []Student{c1, c2} {
			t.Log("Inserted Student:> Id:", c.Id)
			t.Log("Inserted Student:> FirstName:", c.FirstName)
			t.Log("Inserted Student:> LastName:", c.LastName)
			t.Log("Inserted Student:> FullName:", c.FullName)
			t.Log("Inserted Student:> Character:", c.Character)
			t.Log("Inserted Student:> Grade:", c.Grade)
			t.Log("Inserted Student:> Scores:", c.Scores)
			t.Log("Inserted Student:> Description:", c.Description)
			t.Log("")
		}

	}
	t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_One(t *testing.T) {
	t1 := time.Now()
	crew := Student{}
	if e := session.Table(StudentTable).One().Scan(&crew); nil != e {
		t.Fatal("Query One failed:>", e)
	} else {
		t.Log("Queried One:>", crew)
	}
	t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Scan(t *testing.T) {
	t1 := time.Now()
	if r, e := session.Table(StudentTable).Filter(map[string]interface{}{"region": "US"}).All(); nil != e {
		t.Fatal("Table all failed:>", e)
	} else {
		defer r.Close()
		for r.Next() {
			c := Student{}
			if e = r.Scan(&c); nil != e {
				t.Fatal("Scan failed:>", e)
			} else {
				t.Log("Scanned Student:> Id:", c.Id)
				t.Log("Scanned Student:> FirstName:", c.FirstName)
				t.Log("Scanned Student:> LastName:", c.LastName)
				t.Log("Scanned Student:> FullName:", c.FullName)
				t.Log("Scanned Student:> Character:", c.Character)
				t.Log("Scanned Student:> Grade:", c.Grade)
				t.Log("Scanned Student:> Scores:", c.Scores)
				t.Log("Scanned Student:> Description:", c.Description)
			}
		}
	}
	t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Update(t *testing.T) {
	t1 := time.Now()
	var ids []int
	if rows, e := session.Table(StudentTable, "id").LockFor("UPDATE").All(); nil != e {
		t.Fatal("Query with ock failed:>", e)
	} else {
		defer rows.Close()
		for rows.Next() {
			var id int
			if e := rows.Scan(&id); nil != e {
				t.Fatal("Scan failed:>", e)
			} else {
				ids = append(ids, id)
			}
		}
	}
	for _, id := range ids {
		if n, e := session.Table(StudentTable).Where("id", id).Update(map[string]interface{}{"age": 30}); nil != e {
			t.Fatal("Update failed:>", e)
		} else {
			t.Log("Updated rows:>", n)
		}
	}

	t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Delete(t *testing.T) {
	t1 := time.Now()
	if n, e := session.Table(StudentTable).Filter(map[string]interface{}{"full_name": "Tom Cruse"}).Delete(); nil != e {
		t.Fatal("Delete failed:>", e)
	} else {
		t.Log("Deleted rows:>", n)
	}

	if n, e := session.Table(StudentTable).Delete(); nil != e {
		t.Fatal("Delete all failed:>", e)
	} else {
		t.Log("Deleted all rows:>", n)
	}

	t.Log("Time spent:> ", time.Now().Sub(t1))
}

var n int

func Benchmark_Insert(b *testing.B) {
	//var n int
	for i := 0; i < b.N; i++ {
		var c1 = Student{People: People{FullName: fmt.Sprintf("Tom Cruse %d", n), Region: "US", Age: 19 + i, SchoolId: schoolId}}
		c1.FullName = fmt.Sprintf("Tom Cruse %d", n)
		c1.Attributes = make(map[string]interface{})
		c1.Attributes["skill"] = "Good"
		c1.Attributes["score"] = "99.5"
		c1.Scores = []float64{99, 96, 93}
		c1.Character.Attitude = "OK"
		c1.Character.Height = 172
		c1.Character.Weight = 68
		if n, e := session.Table(StudentTable).Insert(&c1); nil != e {
			b.Fatal("Insert failed:> ", e)
		} else {
			b.Log("Inserted Students: ", n)
		}
		n++
	}
}
