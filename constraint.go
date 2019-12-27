package xql

import (
    "strings"
    "fmt"
)

// Constraint Types:
//  - CHECK
//  - NOT NULL
//  - UNIQUE
//  - PRIMARY KEY
//  - FOREIGN KEY
//  - EXCLUDE

const (
    CONSTRAINT_NONE uint8 = iota
    CONSTRAINT_CHECK
    CONSTRAINT_NOT_NULL
    CONSTRAINT_UNIQUE
    CONSTRAINT_PRIMARYKEY
    CONSTRAINT_FOREIGNKEY
    CONSTRAINT_EXCLUDE
    CONSTRAINT_INVALID
)

type Constraint struct {
    Type      uint8
    Columns   []*Column
    Refernces []*Column
    Statement string
    OnDelete  string
    OnUpdate  string
}

func buildConstraints(t *Table, ss... [3]string) []*Constraint {
    var constraints []*Constraint
    for _, xs := range ss {
        constraint := &Constraint{}
        switch strings.ToLower(xs[0]) {
        case "pk", "primarykey":
            constraint.Type = CONSTRAINT_PRIMARYKEY
        case "fk", "foreignkey":
            constraint.Type = CONSTRAINT_FOREIGNKEY
        case "check":
            constraint.Type = CONSTRAINT_CHECK
        case "exclude":
            constraint.Type = CONSTRAINT_EXCLUDE
        case "unique":
            constraint.Type = CONSTRAINT_UNIQUE
        }
        for _,f := range strings.Split(xs[1],",") {
            if field, ok := t.GetColumn(f); ok {
                constraint.Columns = append(constraint.Columns, field)
            }else{
                panic(fmt.Sprintf("Can not get column '%s' from '%s'!", f, t.TableName()))
            }
        }
        constraint.Statement = xs[2]
        constraints = append(constraints, constraint)
    }
    return constraints
}

func makeConstraints(t uint8, fields... interface{}) []*Constraint {
    var constraints []*Constraint
    if t <= CONSTRAINT_NONE || t >= CONSTRAINT_INVALID {
        return constraints
    }
    for _, field := range fields {
        if nil == field {
            continue
        }
        constraint := &Constraint{Type: t}
        if fc, ok := field.(*Column); ok {
            constraint.Columns = []*Column{fc}
            switch t {
            case CONSTRAINT_EXCLUDE:
                if exclude, ok := fc.GetString("exclude"); ok && exclude != "" {
                    constraint.Statement = exclude
                }
            case CONSTRAINT_CHECK:
                if check, ok := fc.GetString("check"); ok && check != "" {
                    constraint.Statement = check
                }else{
                    continue
                }
            case CONSTRAINT_FOREIGNKEY:
                constraint.OnDelete, _ = fc.GetString("ondelete","CASCADE")
                constraint.OnUpdate, _ = fc.GetString("onupdate","CASCADE")
                if fk, ok := fc.GetString("fk"); ok && fk != "" {
                    constraint.Statement = fk
                }else if fk, ok := fc.GetString("foreignkey"); ok && fk != "" {
                    constraint.Statement = fk
                }else{
                    continue
                }
            }
        }else if fcs, ok := field.([]*Column); ok {
            constraint.Columns = fcs
        }else{
            continue
        }
        constraints = append(constraints, constraint)
    }
    return constraints
}