package xql

import (
    "database/sql"
    "fmt"
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
    t.x_columns = make(map[string]*Column)
    t.j_columns = make(map[string]*Column)
    t.m_columns = make(map[string]*Column)
    for _, f := range t.columns {
        t.x_columns[f.FieldName] = f
        t.m_columns[f.ElemName] = f
        t.j_columns[f.Jtag] = f
    }
    if len(schema) > 0 {
        t.schema = schema[0]
    }
    fmt.Println(">>> Table:", t.TableName())
    for _, c := range t.columns {
        fmt.Println(">>> Column:", c.FieldName, c.ElemName, c.TypeDefine)
    }

    for _, c := range t.columns {
        if c.PrimaryKey {
            t.primary_keys = append(t.primary_keys, c)
        }
    }

    //t.constraints = makeConstraints(t.columns...)
    //t.indexes = makeIndexes(t.columns...)
    if tt, ok := entity.(TableConstrainted); ok {
        t.constraints = append(t.constraints, buildConstraints(t, tt.Constraints()...)...)
    }
    if tt, ok := entity.(TableIndexed); ok {
        t.indexes = append(t.indexes, buildIndexes(t, tt.Indexes()...)...)
    }


    return t
}

