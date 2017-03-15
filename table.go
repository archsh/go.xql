package xql

import (
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
        et := reflect.TypeOf(entity)
        for i:=0; i< et.Elem().NumField(); i++ {
            f := et.Elem().Field(i)
            //f.Name
            c := Column{PropertyName:f.Name}
            x_tags := strings.Split(f.Tag.Get("xql"),",")
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
                        t.PrimaryKey = append(t.PrimaryKey, x)
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
            t.Columns[c.FieldName] = c
        }
    }
    return t
}