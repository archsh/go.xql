package xql

import (
    //"fmt"
    "reflect"
    "strings"
)

type Table struct {
    TableName string
    Schema string
    Entity interface{}
    Columns map[string]Column
    PrimaryKey []string

}

type Column struct {
    FieldName string
    PropertyName string
    Type string
    Length uint16
    Unique bool
    Nullable bool
    Indexed bool
    Auto bool
    PrimaryKey bool
}

func _make_columns(entity interface{}) (cols []Column) {
    if nil != entity {
        et := reflect.TypeOf(entity)
        ev := reflect.ValueOf(entity)
        for i:=0; i< et.Elem().NumField(); i++ {
            f := et.Elem().Field(i)
            //f.Name
            x_tags := strings.Split(f.Tag.Get("xql"),",")
            if f.Anonymous  {
                //fmt.Println("_make_columns:>", f.Name, "is Anonymous!", ev.Elem().Field(i))
                if x_tags[0] != "-" {
                    for _, c := range _make_columns(ev.Elem().Field(i).Addr().Interface()) {
                        cols = append(cols, c)
                    }
                }else{
                    continue
                }
                continue
            }
            c := Column{PropertyName:f.Name}
            if len(x_tags) < 1 || x_tags[0]=="" {
                c.FieldName = Camel2Underscore(f.Name)
            }else if x_tags[0] == "-" {
                continue
            }else{
                c.FieldName = x_tags[0]
                for _, x := range x_tags[1:] {
                    switch x {
                    case "pk":
                        c.PrimaryKey = true
                        //t.PrimaryKey = append(t.PrimaryKey, x)
                    case "indexed":
                        c.Indexed = true
                    case "nullable":
                        c.Nullable = true
                    case "unique":
                        c.Unique = true
                    case "auto":
                        c.Auto = true
                    }
                }
            }
            cols = append(cols, c)
            //t.Columns[c.FieldName] = c
        }
    }
    return cols
}

func DeclareTable(name string, entity interface{}, schema ...string) *Table {

    t := &Table{
        TableName:name,
        Entity: entity,
        Columns: make(map[string]Column),
    }
    if len(schema) > 0 {
        t.Schema = schema[0]
    }
    if nil != entity {
        for _, c := range _make_columns(entity) {
            t.Columns[c.FieldName] = c
            if c.PrimaryKey {
                t.PrimaryKey = append(t.PrimaryKey, c.FieldName)
            }
        }
    }
    return t
}