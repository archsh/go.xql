package xql

import (
    "errors"
    "database/sql"
    "fmt"
    log "github.com/Sirupsen/logrus"
)

type Session struct {
    driverName string
    db *sql.DB
    tx *sql.Tx
    verbose bool
}

func (self *Session) getDialect() IDialect {
    if s, ok := _builtin_dialects[self.driverName]; ok {
        return s
    }else{
        panic(fmt.Sprintf("Dialect '%s' not registered! ", self.driverName))
    }
    return nil
}

func (self *Session) Close() {
    if self.tx != nil {
        self.tx.Commit()
    }
}

func (self *Session) Query(table *Table, columns ...interface{}) *QuerySet {
    qs := &QuerySet{session:self, offset:-1, limit:-1}
    qs.table = table
    if len(columns) > 0 {
        for _, c := range columns {
            if qc, ok := c.(QueryColumn); ok {
                qs.queries = append(qs.queries, qc)
            }else if qcn, ok := c.(string); ok {
                if col, ok := qs.table.Columns[qcn]; !ok {
                    panic("Invalid column name:"+qcn)
                }else{
                    qs.queries = append(qs.queries, QueryColumn{FieldName:col.FieldName, Alias:col.FieldName})
                }
            }else{
                panic("Unsupported parameter type!")
            }

        }
    }else{
        //for _, col := range qs.table.Columns {
        //    qs.queries = append(qs.queries, QueryColumn{FieldName:col.FieldName, Alias:col.FieldName})
        //}
    }
    return qs
}

func (self *Session) Begin() error {
    if nil != self.tx {
        return errors.New("Already in Tx!")
    }
    tx, err := self.db.Begin()
    if nil == err {
        self.tx = tx
    }
    return err
}

func (self *Session) Commit() error {
    if nil == self.tx {
        return errors.New("Not open Tx!")
    }
    err := self.tx.Commit()
    if nil == err {
        self.tx = nil
    }
    return err
}

func (self *Session) Rollback() error {
    if nil == self.tx {
        return errors.New("Not open Tx!")
    }
    err := self.tx.Rollback()
    self.tx = nil
    return err
}

func (self *Session) doExec(query string, args ...interface{}) (sql.Result, error) {
    if self.tx != nil {
        if self.verbose { log.Debugln("doExec in Tx: ", query, args) }
        return self.tx.Exec(query, args...)
    }else{
        if self.verbose { log.Debugln("doExec in DB: ", query, args) }
        return self.db.Exec(query, args...)
    }
    return nil, nil
}

func (self *Session) doQuery(query string, args ...interface{}) (*sql.Rows, error) {
    if self.tx != nil {
        if self.verbose { log.Debugln("doQuery in Tx: ", query, args) }
        return self.tx.Query(query, args...)
    }else{
        if self.verbose { log.Debugln("doQuery in DB: ", query, args) }
        return self.db.Query(query, args...)
    }
    return nil, nil
}

func (self *Session) doQueryRow(query string, args ...interface{}) *sql.Row {
    if self.tx != nil {
        if self.verbose { log.Debugln("doQueryRow in Tx: ", query, args) }
        return self.tx.QueryRow(query, args...)
    }else{
        if self.verbose { log.Debugln("doQueryRow in Db: ", query, args) }
        return self.db.QueryRow(query, args...)
    }
    return nil
}