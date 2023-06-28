package xql

import (
	"fmt"
	"strings"
)

// Constraint Types:
//  - CHECK
//  - NOT NULL
//  - UNIQUE
//  - PRIMARY KEY
//  - FOREIGN KEY
//  - EXCLUDE

const (
	ConstraintNone uint8 = iota
	ConstraintCheck
	ConstraintNotNull
	ConstraintUnique
	ConstraintPrimaryKey
	ConstraintForeignKey
	ConstraintExclude
	ConstraintInvalid
)

type Constraint struct {
	Type       uint8
	Columns    []*Column
	References []*Column
	Statement  string
	OnDelete   string
	OnUpdate   string
}

func buildConstraints(t *Table, ss ...[3]string) []*Constraint {
	var constraints []*Constraint
	for _, xs := range ss {
		constraint := &Constraint{}
		switch strings.ToLower(xs[0]) {
		case "pk", "primarykey":
			constraint.Type = ConstraintPrimaryKey
		case "fk", "foreignkey":
			constraint.Type = ConstraintForeignKey
		case "check":
			constraint.Type = ConstraintCheck
		case "exclude":
			constraint.Type = ConstraintExclude
		case "unique":
			constraint.Type = ConstraintUnique
		}
		for _, f := range strings.Split(xs[1], ",") {
			if field, ok := t.GetColumn(f); ok {
				constraint.Columns = append(constraint.Columns, field)
			} else {
				panic(fmt.Sprintf("Can not get column '%s' from '%s'!", f, t.TableName()))
			}
		}
		constraint.Statement = xs[2]
		constraints = append(constraints, constraint)
	}
	return constraints
}

func makeConstraints(t uint8, fields ...interface{}) []*Constraint {
	var constraints []*Constraint
	if t <= ConstraintNone || t >= ConstraintInvalid {
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
			case ConstraintExclude:
				if exclude, ok := fc.GetString("exclude"); ok && exclude != "" {
					constraint.Statement = exclude
				}
			case ConstraintCheck:
				if check, ok := fc.GetString("check"); ok && check != "" {
					constraint.Statement = check
				} else {
					continue
				}
			case ConstraintForeignKey:
				constraint.OnDelete, _ = fc.GetString("ondelete", "CASCADE")
				constraint.OnUpdate, _ = fc.GetString("onupdate", "CASCADE")
				if fk, ok := fc.GetString("fk"); ok && fk != "" {
					constraint.Statement = fk
				} else if fk, ok := fc.GetString("foreignkey"); ok && fk != "" {
					constraint.Statement = fk
				} else {
					continue
				}
			}
		} else if fcs, ok := field.([]*Column); ok {
			constraint.Columns = fcs
		} else {
			continue
		}
		constraints = append(constraints, constraint)
	}
	return constraints
}
