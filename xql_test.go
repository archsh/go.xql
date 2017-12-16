package xql

import (
    "testing"
    "time"
    _ "github.com/lib/pq"
    "github.com/archsh/go.uuid"
    _ "github.com/archsh/go.xql/dialects/postgres"

)

type Crew struct {
    Id          string `json:"id"`
    FullName    string `json:"fullName" xql:"fullname"`
    FirstName   string `json:"firstName"`
    MiddleName  string `json:"middleName"`
    LastName    string `json:"lastName"`
    Region      string `json:"region"`
    ImdbId      string `json:"imdbId"`
    Description string `json:"description"`
    Created     *time.Time `json:"created"`
    Updated     *time.Time `json:"Updated"`
}

var MovieCrew = DeclareTable("metas_crews", &Crew{}, "deneb")

func TestCreateEngine(t *testing.T) {
    t1 := time.Now()
    engine, e := CreateEngine("postgres",
        "host=localhost port=5432 user=postgres password=postgres dbname=cygnuxdb sslmode=disable")
    if nil != e {
        t.Fatal("Connec Database failed:> ", e)
    }
    t.Log("MovieCrew:> ", MovieCrew)
    _ = engine.MakeSession()
    t.Log("Time spent:> ", time.Now().Sub(t1))
}

func TestQuerySet_Insert(t *testing.T) {
    t1 := time.Now()
    engine, e := CreateEngine("postgres",
        "host=localhost port=5432 user=postgres password=postgres dbname=cygnuxdb sslmode=disable")
    if nil != e {
        t.Fatal("Connec Database failed:> ", e)
    }
    t.Log("MovieCrew:> ", MovieCrew)
    sess := engine.MakeSession()
    c1 := Crew{Id:uuid.NewV4().String(), FullName:"Tom Cruse", Region:"US"}
    c2 := Crew{Id:uuid.NewV4().String(), FullName:"Hue Jackman", Region:"US"}
    n, e := sess.Query(MovieCrew).Insert(&c1, &c2)
    if nil != e {
        t.Fatal("Insert failed:> ", e)
    }
    t.Log("Insert lines:>", n)
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
