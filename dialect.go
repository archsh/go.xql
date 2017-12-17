package xql

import (
    "database/sql"
)

type IDialect interface {
    Create(*Table, ...interface{}) (string, []interface{}, error)
    Select(*Table, []QueryColumn, []QueryFilter, []QueryOrder, int64, int64) (string, []interface{}, error)
    Insert(*Table, interface{}, ...string) (string, []interface{}, error)
    Update(*Table, []QueryFilter, ...UpdateColumn) (string, []interface{}, error)
    Delete(*Table, []QueryFilter) (string, []interface{}, error)
}

var _builtin_dialects map[string]IDialect

func init() {
    _builtin_dialects = make(map[string]IDialect)
}

func RegisterDialect(name string, d IDialect) {
    _builtin_dialects[name] = d
}

type Engine struct {
    db         *sql.DB
    driverName string
}

func CreateEngine(name string, dataSource string) (*Engine, error) {
    db, err := sql.Open(name, dataSource)
    if nil != err {
        return nil, err
    }
    return &Engine{db: db, driverName: name}, nil
}

func (engine *Engine) MakeSession() *Session {
    return &Session{
        db:         engine.db,
        driverName: engine.driverName,
    }
}
