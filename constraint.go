package xql

// Constraint Types:
//  - Check
//  - Not-Null
//  - UNIQUE
//  - PRIMARY KEY
//  - FOREIGN KEY
//  - Exclusion

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


func makeConstraints(fields... interface{}) []*Constraint {
    var constraints []*Constraint
    for _, field := range fields {
        if nil == field {
            continue
        }
    }
    return constraints
}