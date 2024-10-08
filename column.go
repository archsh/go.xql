package xql

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

type Column struct {
	PropertySet
	FieldName   string
	ElemName    string
	JTag        string
	Type        reflect.Type
	TypeDefine  string
	Indexed     bool // Indexed or not, on field
	Nullable    bool // Nullable constraint on field
	Unique      bool // Unique constraint on field
	PrimaryKey  bool //Primary Key constraint on field
	Default     interface{}
	Constraints []*Constraint
	Indexes     []*Index
	table       interface{}
}

type Declarable interface {
	Declare(props PropertySet) string
}

func DefaultDeclare(f reflect.StructField, props PropertySet) (string, error) {
	if t, ok := props.GetString("type"); ok {
		t = strings.ToLower(t)
		switch t {
		case "varchar", "string":
			return Varchar("").Declare(props), nil
		case "char":
			return Char("").Declare(props), nil
		case "text":
			return Text("").Declare(props), nil
		case "int", "integer":
			return Integer(0).Declare(props), nil
		case "smallint", "smallinteger":
			return SmallInteger(0).Declare(props), nil
		case "bigint", "biginteger":
			return BigInteger(0).Declare(props), nil
		case "serial":
			return Serial(0).Declare(props), nil
		case "bigserial":
			return BigSerial(0).Declare(props), nil
		case "real", "float":
			return Real(0.0).Declare(props), nil
		case "double":
			return Double(0.0).Declare(props), nil
		case "bool", "boolean":
			return Boolean(false).Declare(props), nil
		case "date":
			return Date(time.Time{}).Declare(props), nil
		case "time":
			return Time(time.Time{}).Declare(props), nil
		case "datetime", "timestamp":
			return TimeStamp(time.Time{}).Declare(props), nil
		case "decimal", "numeric":
			return Decimal("").Declare(props), nil
		case "uuid":
			return UUID("").Declare(props), nil
		default:
			return t, nil
		}
	}
	switch f.Type.Kind() {
	case reflect.String:
		//size, _ := props.GetUInt("size", 32)
		//return fmt.Sprintf("VARCHAR(%d)", size), nil
		return Varchar("").Declare(props), nil
	case reflect.Int16, reflect.Uint16:
		//return "SMALLINT", nil
		return SmallInteger(0).Declare(props), nil
	case reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32:
		//return "INTEGER", nil
		return Integer(0).Declare(props), nil
	case reflect.Int64, reflect.Uint64:
		//return "BIGINT", nil
		return BigInteger(0).Declare(props), nil
	case reflect.Bool:
		return Boolean(false).Declare(props), nil
	case reflect.Float32:
		//return "FLOAT", nil
		return Real(0.0).Declare(props), nil
	case reflect.Float64:
		return Double(0.0).Declare(props), nil
	default:
		return "", errors.New("Unknown type of:>" + f.Name)
	}

}

// makeColumn
// Make a &Column{} object according to given field.
func makeColumn(t *Table, f reflect.StructField, v reflect.Value) *Column {
	props, e := ParseProperties(f.Tag.Get("xql"))
	if nil != e {
		panic(e)
	}

	field := &Column{
		FieldName:   Camel2Underscore(f.Name),
		ElemName:    f.Name,
		Type:        f.Type,
		PropertySet: props,
	}
	jtag := f.Tag.Get("json")
	if jtag != "" {
		field.JTag = strings.Split(jtag, ",")[0]
	} else {
		field.JTag = f.Name
	}
	if fn, ok := props.PopString("name"); ok {
		field.FieldName = fn
	}
	field.Indexed, _ = props.PopBool("index", false)
	if field.Indexed {
		field.Indexes = append(field.Indexes,
			makeIndexes(IndexBTree, t.BaseTableName()+"_"+field.FieldName, field)...)
	}
	field.Nullable, _ = props.PopBool("nullable", false)
	if field.Nullable == false {
		field.Constraints = append(field.Constraints,
			makeConstraints(ConstraintNotNull, field)...)
	}
	field.Unique, _ = props.PopBool("unique", false)
	if field.Unique {
		field.Constraints = append(field.Constraints,
			makeConstraints(ConstraintUnique, field)...)
	}
	field.PrimaryKey, _ = props.PopBool("primarykey", false)
	if !field.PrimaryKey {
		field.PrimaryKey, _ = props.PopBool("pk", false)
	}
	if field.PrimaryKey {
		field.Constraints = append(field.Constraints,
			makeConstraints(ConstraintPrimaryKey, field)...)
	}
	if fk, ok := props.GetString("foreignkey"); ok && fk != "" {
		field.Constraints = append(field.Constraints,
			makeConstraints(ConstraintForeignKey, field)...)
	} else if fk, ok := props.GetString("fk"); ok && fk != "" {
		field.Constraints = append(field.Constraints,
			makeConstraints(ConstraintForeignKey, field)...)
	}
	if check, ok := props.GetString("check"); ok && check != "" {
		field.Constraints = append(field.Constraints,
			makeConstraints(ConstraintCheck, field)...)
	}
	if exclude, ok := props.GetString("exclude"); ok && exclude != "" {
		field.Constraints = append(field.Constraints,
			makeConstraints(ConstraintExclude, field)...)
	}
	if df, ok := props.PopString("default"); ok {
		field.Default = df
	}
	//field.PropertySet = props
	if p, ok := v.Interface().(Declarable); ok {
		field.TypeDefine = p.Declare(props)
	} else {
		if d, e := DefaultDeclare(f, props); nil == e {
			field.TypeDefine = d
		} else {
			//fmt.Println("makeColumn:>", props)
			panic(e)
		}
	}
	return field
}

// makeColumns
// Make a list of &Column{} objects according to a given struct pointer.
func makeColumns(t *Table, p interface{}, recursive bool, skips ...string) []*Column {
	if nil == p {
		panic("Can not use nil pointer ")
	}
	var et reflect.Type
	var ev reflect.Value
	tt := reflect.TypeOf(p)
	vv := reflect.ValueOf(p)
	if tt.Kind() == reflect.Interface || tt.Kind() == reflect.Ptr || tt.Kind() == reflect.UnsafePointer {
		et = tt.Elem()
		ev = vv.Elem()

	} else {
		et = tt
		ev = vv
	}
	var fields []*Column
	for i := 0; i < et.NumField(); i++ {
		f := et.Field(i)
		v := ev.Field(i)
		if inSlice(f.Name, skips) {
			continue
		}
		xTags := strings.Split(f.Tag.Get("xql"), ",")
		if xTags[0] == "-" {
			continue
		}
		if f.Anonymous {
			if xTags[0] != "-" {
				sks := getSkips(xTags)
				for _, c := range makeColumns(t, ev.Field(i).Interface(), true, sks...) {
					fields = append(fields, c)
				}
			} else {
				continue
			}
			continue
		}
		field := makeColumn(t, f, v)
		fields = append(fields, field)
	}
	return fields
}
