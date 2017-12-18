package xql

// Constraint Types:
//  - Check
//  - Not-Null
//  - UNIQUE
//  - PRIMARY KEY
//  - FOREIGN KEY
//  - Exclusion

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
    Type  uint8
    Fields []*Column
    Refernces []*Column
}

type ForeignKey struct {
    Constraint
}

type Unique struct {
    Constraint
}

type PrimaryKey struct {
    Constraint
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
            constraint.Fields = []*Column{fc}
        }else if fcs, ok := field.([]*Column); ok {
            constraint.Fields = fcs
        }else{
            continue
        }
        constraints = append(constraints, constraint)
    }
    return constraints
}