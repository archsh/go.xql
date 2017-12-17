package xql

import (
    "reflect"
    "strings"
    //"fmt"
)

// Table ...
// Struct defined for a table object.
type Table struct {
    fields       []*Field
    constraints  []*Constraint
    primary_keys []*Field
    foreign_keys []*Field
}

type TableIdentified interface {
    TableName() string
    TableFullName() string
    Construct() error
}

type TableSQLed interface {
    Create(session *Session) error
    Drop(session *Session) error
    Truncate(session *Session) error
    Select(set QuerySet)
    Update(set QuerySet)
    Insert()
    Delete()
}

func (t Table) TableName() string {
    return ""
}

func (t Table) TableFullName() string {
    return t.TableName()
}

func (t Table) Construct() error {
    return nil
}

func GetFields(t interface{}) ([]*Field, error) {
    fields := []*Field{}

    return fields, nil
}

func GetConstraints(t interface{}) ([]*Constraint, error) {
    constraints := []*Constraint{}

    return constraints, nil
}

func GetIndexes(t interface{}) ([]*Index, error) {
    indexes := []*Index{}

    return indexes, nil
}

func inSlice(a string, ls []string) bool {
    for _, s := range ls {
        if a == s {
            return true
        }
    }
    return false
}

func getSkips(tags []string) (skips []string) {
    if nil == tags || len(tags) < 1 {
        return
    }
    for _, tag := range tags {
        if strings.HasPrefix(tag, "skips:") {
            s := strings.TrimLeft(tag, "skips:")
            for _, n := range strings.Split(s, ";") {
                if n != "" {
                    skips = append(skips, n)
                }
            }
            return
        }
    }
    return
}

func makeColumns(entity interface{}, recursive bool, skips ...string) (cols []*Column) {
    if nil != entity {
        et := reflect.TypeOf(entity)
        ev := reflect.ValueOf(entity)
        for i := 0; i < et.Elem().NumField(); i++ {
            f := et.Elem().Field(i)
            //fv := et.Elem().Field(i)
            if inSlice(f.Name, skips) {
                continue
            }
            x_tags := strings.Split(f.Tag.Get("xql"), ",")
            if f.Anonymous && !recursive {
                if x_tags[0] != "-" {
                    sks := getSkips(x_tags)
                    for _, c := range makeColumns(ev.Elem().Field(i).Addr().Interface(), true, sks...) {
                        cols = append(cols, c)
                    }
                } else {
                    continue
                }
                continue
            }
            c := &Column{PropertyName: f.Name}
            if len(x_tags) < 1 {
                c.FieldName = Camel2Underscore(f.Name)
            } else if x_tags[0] == "-" {
                continue
            }
            if x_tags[0] == "" {
                c.FieldName = Camel2Underscore(f.Name)
            } else {
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
                    default:
                        xs := strings.Split(x, "=")
                        if len(xs) > 1 {
                            switch xs[0] {
                            case "type":
                                c.Type = xs[1]
                            }
                        }
                    }
                }
            }
            json_tags := strings.Split(f.Tag.Get("json"), ",")
            if len(json_tags) < 1 {
                c.JTAG = c.PropertyName
            } else if json_tags[0] != "-" {
                if json_tags[0] == "" {
                    c.JTAG = c.PropertyName
                } else {
                    c.JTAG = json_tags[0]
                }
            }
            if c.Type == "" {
                switch f.Type.Kind() {
                case reflect.String:
                    c.Type = "VARCHAR(32)"
                case reflect.Int16, reflect.Uint16:
                    c.Type = "SMALLINT"
                case reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32:
                    c.Type = "INTEGER"
                case reflect.Int64, reflect.Uint64:
                    c.Type = "BIGINT"
                case reflect.Bool:
                    c.Type = "BOOLEAN"
                case reflect.Float32, reflect.Float64:
                    c.Type = "FLOAT"
                }
            }
            cols = append(cols, c)
            //t.MappedColumns[c.FieldName] = c
        }
    }
    return cols
}

func DeclareTable(name string, entity interface{}, schema ...string) *Table {
    t := &Table{
        TableName:      name,
        Entity:         entity,
        MappedColumns:  make(map[string]*Column),
        JTaggedColumns: make(map[string]*Column),
    }
    if len(schema) > 0 {
        t.Schema = schema[0]
    }
    if nil != entity {
        for _, c := range makeColumns(entity, false) {
            t.MappedColumns[c.FieldName] = c
            t.ListedColumns = append(t.ListedColumns, c)
            if c.JTAG != "" {
                t.JTaggedColumns[c.JTAG] = c
            }
            if c.PrimaryKey {
                t.PrimaryKey = append(t.PrimaryKey, c.FieldName)
            }
        }
    }
    return t
}

func (t *Table) GetColumn(name string) (*Column, bool) {
    if c, ok := t.MappedColumns[name]; ok {
        return c, true
    }
    if c, ok := t.JTaggedColumns[name]; ok {
        return c, true
    }
    return nil, false
}
