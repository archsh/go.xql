package xql

import (
    "errors"
    "database/sql"
    "fmt"
    "time"
)

type Session struct {
    driverName string
    dialect    IDialect
    db         *sql.DB
    tx         *sql.Tx
    verbose    bool
}

func (self *Session) getDialect() IDialect {
    if nil != self.dialect {
        return self.dialect
    }
    if s, ok := _builtin_dialects[self.driverName]; ok {
        self.dialect = s
        return s
    } else {
        panic(fmt.Sprintf("Dialect '%s' not registered! ", self.driverName))
    }
    return nil
}

func (self *Session) Drop(table *Table, force bool) error {
    s, args, e := self.getDialect().Drop(table, force)
    if nil != e {
        return e
    }
    if _, e := self.db.Exec(s, args...); nil != e {
        return errors.New(e.Error()+":>"+s)
    } else {
        return nil
    }
}

func (self *Session) Create(table *Table) error {
    s, args, e := self.getDialect().Create(table)
    if nil != e {
        return e
    }
    //log.Debugln("SQL:>>>", s)
    if _, e := self.db.Exec(s, args...); nil != e {
        return errors.New(e.Error()+":>"+s)
    } else {
        return nil
    }
    //log.Debugln("SQL:>>>", s)
    //return e
}

func (self *Session) Exec(s string,args...interface{}) (sql.Result, error) {
    return self.doExec(s, args...)
    //if _, e := self.db.Exec(s); nil != e {
    //    return e
    //} else {
    //    return nil
    //}
}

func (self *Session) Close() {
    if self.tx != nil {
        self.tx.Rollback()
    }
}

func (self *Session) Query(table *Table, columns ...interface{}) QuerySet {
    qs := QuerySet{session: self, offset: -1, limit: -1}
    qs.table = table
    if len(columns) > 0 {
        for i, c := range columns {
            if qc, ok := c.(QueryColumn); ok {
                qs.queries = append(qs.queries, qc)
            } else if qcn, ok := c.(string); ok {
                if col, ok := qs.table.GetColumn(qcn); !ok {
                    //panic("Invalid column name:" + qcn)
                    qs.queries = append(qs.queries, QueryColumn{FieldName: qcn, Alias: fmt.Sprintf("aa%d",i)})
                } else {
                    qs.queries = append(qs.queries, QueryColumn{FieldName: col.FieldName, Alias: col.FieldName})
                }
            } else {
                panic("Unsupported parameter type!")
            }

        }
    } else {
        //for _, col := range qs.table.MappedColumns {
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

func log_timing(t1 time.Time, msg string, params ...interface{}) {
    t2 := time.Now()
    fmt.Println(t1, t2, t2.Sub(t1), msg, params)
}

func (self *Session) doExec(query string, args ...interface{}) (sql.Result, error) {
    if self.verbose {
        t1 := time.Now()
        defer log_timing(t1,"Session.doExec:", query)
    }
    if self.tx != nil {
        if self.verbose {
            //log.Debugln("doExec in Tx: ", query, args)
        }
        return self.tx.Exec(query, args...)
    } else {
        if self.verbose {
            //log.Debugln("doExec in DB: ", query, args)
        }
        return self.db.Exec(query, args...)
    }
    return nil, nil
}

func (self *Session) doQuery(query string, args ...interface{}) (*sql.Rows, error) {
    if self.verbose {
        t1 := time.Now()
        defer log_timing(t1,"Session.doQuery:", query)
    }
    if self.tx != nil {
        if self.verbose {
            //log.Debugln("doQuery in Tx: ", query, args)
        }
        return self.tx.Query(query, args...)
    } else {
        if self.verbose {
            //log.Debugln("doQuery in DB: ", query, args)
        }
        return self.db.Query(query, args...)
    }
    return nil, nil
}

func (self *Session) doQueryRow(query string, args ...interface{}) *sql.Row {
    if self.verbose {
        t1 := time.Now()
        defer log_timing(t1,"Session.doQueryRaw:", query)
    }
    if self.tx != nil {
        if self.verbose {
            //log.Debugln("doQueryRow in Tx: ", query, args)
        }
        return self.tx.QueryRow(query, args...)
    } else {
        if self.verbose {
            //log.Debugln("doQueryRow in Db: ", query, args)
        }
        return self.db.QueryRow(query, args...)
    }
    return nil
}
