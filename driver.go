package xql

import (
    "database/sql"
)

type StatementBuilder interface {
    Create(*Table, ...interface{}) (string, []interface{}, error)
    Select(*Table, []QueryColumn, []QueryFilter, []QueryOrder, int64, int64) (string, []interface{}, error)
    Insert(*Table, interface{}, ...string) (string, []interface{}, error)
    Update(*Table, []QueryFilter, ...UpdateColumn) (string, []interface{}, error)
    Delete(*Table, []QueryFilter) (string, []interface{}, error)
}


var _statement_builders map[string]StatementBuilder

func init() {
    _statement_builders = make(map[string]StatementBuilder)
}

func RegisterBuilder(name string, d StatementBuilder) {
    _statement_builders[name] = d
}

type Engine struct {
    db *sql.DB
    driverName string
}

func CreateEngine(name string, dataSource string) (*Engine, error) {
    db, err := sql.Open(name, dataSource)
    if nil != err {
        return nil, err
    }
    return &Engine{db:db, driverName:name}, nil
}

func (engine *Engine) MakeSession() *Session {
    return &Session{
        db: engine.db,
        driverName: engine.driverName,
    }
}