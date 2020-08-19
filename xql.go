package xql

import (
    "database/sql"
    //"fmt"
)

func MakeSession(db *sql.DB, driverName string, verbose ...bool) *Session {
    sess := &Session{db: db, driverName: driverName}
    if len(verbose) > 0 {
        sess.verbose = verbose[0]
    }
    return sess
}

// DeclareTable
// Which declare a new Table instance according to a given entity.
func DeclareTable(entity TableIdentified, schema ...string) *Table {
    var skips []string
    if et, ok := entity.(TableIgnored); ok {
        skips = et.Ignore()
    }
    t := &Table{
        entity:  entity,
    }
    t.columns = makeColumns(t, entity, false, skips...)
    t.xColumns = make(map[string]*Column)
    t.jColumns = make(map[string]*Column)
    t.mColumns = make(map[string]*Column)
    for _, f := range t.columns {
        t.xColumns[f.FieldName] = f
        t.mColumns[f.ElemName] = f
        t.jColumns[f.JTag] = f
    }
    if len(schema) > 0 {
        t.schema = schema[0]
    }
    //fmt.Println(">>> Table:", t.TableName())
    //for _, c := range t.columns {
    //    fmt.Println(">>> Column:", c.FieldName, c.ElemName, c.TypeDefine)
    //}

    //t.constraints = makeConstraints(t.columns...)
    //t.indexes = makeIndexes(t.columns...)
    if tt, ok := entity.(TableConstrained); ok {
        t.constraints = append(t.constraints, buildConstraints(t, tt.Constraints()...)...)
    }
    if tt, ok := entity.(TableIndexed); ok {
        t.indexes = append(t.indexes, buildIndexes(t, tt.Indexes()...)...)
    }

    for _, c := range t.constraints {
        if c.Type == ConstraintPrimaryKey {
            t.primaryKeys = append(t.primaryKeys, c.Columns...)
        }
    }

    for _, c := range t.columns {
        if c.PrimaryKey {
            t.primaryKeys = append(t.primaryKeys, c)
        }
    }

    return t
}

