package xql

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
                constraint.OnDelete, _ = fc.GetString("ondelete")
                constraint.OnUpdate, _ = fc.GetString("onupdate")
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