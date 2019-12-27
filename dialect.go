package xql

import (
	"database/sql"
)

type IDialect interface {
	Create(*Table, ...interface{}) (string, []interface{}, error)
	Drop(*Table, bool) (string, []interface{}, error)
	Select(*Table, []QueryColumn, []QueryFilter, []QueryOrder, QueryExtra, int64, int64) (string, []interface{}, error)
	Insert(*Table, interface{}, ...string) (string, []interface{}, error)
	InsertWithInsertedId(*Table, interface{}, string, ...string) (string, []interface{}, error)
	Update(*Table, []QueryFilter, QueryExtra, ...UpdateColumn) (string, []interface{}, error)
	Delete(*Table, []QueryFilter, QueryExtra) (string, []interface{}, error)
}

var builtinDialects map[string]IDialect

func init() {
	builtinDialects = make(map[string]IDialect)
}

func RegisterDialect(name string, d IDialect) {
	builtinDialects[name] = d
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

func (engine Engine) DB() *sql.DB {
	return engine.db
}

func (engine Engine) DriverName() string {
	return engine.driverName
}

func (engine *Engine) MakeSession() *Session {
	return &Session{
		db:         engine.db,
		driverName: engine.driverName,
	}
}
