package sqlite

import (
	"fmt"
	xql "github.com/archsh/go.xql"
)

type sqliteDialect struct{}

func (s sqliteDialect) Create(table *xql.Table, i ...interface{}) (stm string, args []interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteDialect) Drop(table *xql.Table, b bool) (stm string, args []interface{}, err error) {
	//TODO implement me
	//panic("implement me")
	stm = fmt.Sprintf("DROP TABLE %s", table.TableName())
	return
}

func (s sqliteDialect) Select(table *xql.Table, columns []xql.QueryColumn, filters []xql.QueryFilter, orders []xql.QueryOrder, s2 string, i int64, i2 int64) (stm string, args []interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteDialect) Insert(table *xql.Table, i interface{}, s2 ...string) (stm string, args []interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteDialect) InsertWithInsertedId(table *xql.Table, i interface{}, s3 string, s2 ...string) (stm string, args []interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteDialect) Update(table *xql.Table, filters []xql.QueryFilter, column ...xql.UpdateColumn) (stm string, args []interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (s sqliteDialect) Delete(table *xql.Table, filters []xql.QueryFilter) (stm string, args []interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	xql.RegisterDialect("sqlite", &sqliteDialect{})
}
