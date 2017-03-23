package xql

import (
    "reflect"
    "strings"
    //"fmt"
)

type Table struct {
    TableName string
    Schema string
    Entity interface{}
    Columns map[string]*Column
    PrimaryKey []string
    columns_by_jtag map[string]*Column

}

type Column struct {
    FieldName    string
    PropertyName string
    JTAG         string
    Type         string
    Length       uint16
    Unique       bool
    Nullable     bool
    Indexed      bool
    Auto         bool
    PrimaryKey   bool
}

func _in_slice(a string, ls []string) bool {
    for _, s := range ls {
        if a == s {
            return true
        }
    }
    return false
}

func _get_skips(tags []string) (skips []string) {
    if nil == tags || len(tags) < 1 {
        return
    }
    for _, tag := range tags {
        if strings.HasPrefix(tag, "skips:") {
            s := strings.TrimLeft(tag, "skips:")
            for _, n := range strings.Split(s, ";"){
                if n != "" {
                    skips = append(skips, n)
                }
            }
            return
        }
    }
    return
}

func _make_columns(entity interface{}, recursive bool, skips ...string) (cols []*Column) {
    if nil != entity {
        et := reflect.TypeOf(entity)
        ev := reflect.ValueOf(entity)
        for i:=0; i< et.Elem().NumField(); i++ {
            f := et.Elem().Field(i)
            if _in_slice(f.Name, skips) {
                continue
            }
            x_tags := strings.Split(f.Tag.Get("xql"),",")
            if f.Anonymous && ! recursive {
                if x_tags[0] != "-" {
                    sks := _get_skips(x_tags)
                    for _, c := range _make_columns(ev.Elem().Field(i).Addr().Interface(),true, sks...) {
                        cols = append(cols, c)
                    }
                }else{
                    continue
                }
                continue
            }
            c := &Column{PropertyName:f.Name}
            if len(x_tags) < 1 {
                c.FieldName = Camel2Underscore(f.Name)
            }else if x_tags[0] == "-" {
                continue
            }
            if x_tags[0]=="" {
                c.FieldName = Camel2Underscore(f.Name)
            }else{
                c.FieldName = x_tags[0]
            }
            if len(x_tags) > 1 {
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
            json_tags := strings.Split(f.Tag.Get("json"),",")
            if len(json_tags) < 1 {
                c.JTAG = c.PropertyName
            }else if json_tags[0] != "-" {
                if json_tags[0] == "" {
                    c.JTAG = c.PropertyName
                }else{
                    c.JTAG = json_tags[0]
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
        Columns: make(map[string]*Column),
        columns_by_jtag: make(map[string]*Column),
    }
    if len(schema) > 0 {
        t.Schema = schema[0]
    }
    if nil != entity {
        for _, c := range _make_columns(entity, false) {
            t.Columns[c.FieldName] = c
            if c.JTAG != "" {
                t.columns_by_jtag[c.JTAG] = c
            }
            if c.PrimaryKey {
                t.PrimaryKey = append(t.PrimaryKey, c.FieldName)
            }
        }
    }
    return t
}

func (t *Table) GetColumn(name string) (*Column, bool) {
    if c, ok := t.Columns[name]; ok {
        return c, true
    }
    if c, ok := t.columns_by_jtag[name]; ok {
        return c, true
    }
    return nil, false
}