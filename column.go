package xql

import (
    "reflect"
    "strings"
    "errors"
    "fmt"
)

type Column struct {
    PropertySet
    FieldName  string
    ElemName   string
    Jtag       string
    Type       reflect.Type
    TypeDefine string
    Indexed    bool // Indexed or not, on field
    Nullable   bool // Nullable constraint on field
    Unique     bool // Unique constraint on field
    Check      interface{}  // Check constraint on field
    PrimaryKey bool //Primary Key constraint on field
    Default    interface{}
    ForeignKey interface{}
    table      interface{}
}

type Declarable interface {
    Declare(props PropertySet) string
}

func GenericDeclare(f reflect.StructField, props PropertySet) (string, error) {
    switch f.Type.Kind() {
    case reflect.String:
        size, _ := props.GetUInt("size", 32)
        return fmt.Sprintf("VARCHAR(%d)", size), nil
    case reflect.Int16, reflect.Uint16:
        return "SMALLINT", nil
    case reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32:
        return "INTEGER", nil
    case reflect.Int64, reflect.Uint64:
        return "BIGINT", nil
    case reflect.Bool:
        return "BOOLEAN", nil
    case reflect.Float32, reflect.Float64:
        return "FLOAT", nil
    }
    return "", errors.New("Unknow Type!")
}

// makeColumn
// Make a &Column{} object according to given field.
func makeColumn(f reflect.StructField, v reflect.Value) *Column {
    props, e := ParseProperties(f.Tag.Get("xql"))
    if nil != e {
        panic(e)
    }
    field := &Column{
        FieldName:Camel2Underscore(f.Name),
        ElemName:f.Name,
        Type: f.Type,
    }
    jtag := f.Tag.Get("json")
    if jtag != "" {
        field.Jtag = jtag
    }else{
        field.Jtag = f.Name
    }
    field.Indexed, _ = props.PopBool("index", false)
    field.Nullable, _ = props.PopBool("nullable", false)
    field.Unique, _ = props.PopBool("unique", false)
    if fn, ok := props.PopString("name"); ok {
        field.FieldName = fn
    }
    if df, ok := props.PopString("default"); ok {
        field.Default = df
    }
    field.PropertySet = props
    if p, ok := v.Interface().(Declarable); ok {
        field.TypeDefine = p.Declare(props)
    }else{
        if d, e := GenericDeclare(f, props); nil == e {
            field.TypeDefine = d
        }else{
            panic(e)
        }
    }
    return field
}

// makeColumns
// Make a list of &Column{} objects according to a given struct pointer.
func makeColumns(p interface{}, recursive bool, skips ...string) []*Column {
    if nil == p {
        panic("Can not use nil pointer ")
    }

    et := reflect.TypeOf(p)
    ev := reflect.ValueOf(p)
    fields := make([]*Column,et.Elem().NumField())
    for i := 0; i < et.Elem().NumField(); i++ {
        f := et.Elem().Field(i)
        v := ev.Elem().Field(i)
        if inSlice(f.Name, skips) {
            continue
        }
        x_tags := strings.Split(f.Tag.Get("xql"), ",")
        if f.Anonymous && !recursive {
            if x_tags[0] != "-" {
                sks := getSkips(x_tags)
                for _, c := range makeColumns(ev.Elem().Field(i).Addr().Interface(), true, sks...) {
                    fields = append(fields, c)
                }
            } else {
                continue
            }
            continue
        }
        field := makeColumn(f, v)
        fields = append(fields, field)
    }
    return fields
}


