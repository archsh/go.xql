package xql

import (
	"database/sql"
	"errors"
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

func (session *Session) getDialect() IDialect {
	if nil != session.dialect {
		return session.dialect
	}
	if s, ok := builtinDialects[session.driverName]; ok {
		session.dialect = s
		return s
	} else {
		panic(fmt.Sprintf("Dialect '%s' not registered! ", session.driverName))
	}
	return nil
}

func (session *Session) Drop(table *Table, force bool) error {
	s, args, e := session.getDialect().Drop(table, force)
	if nil != e {
		return e
	}
	if _, e := session.db.Exec(s, args...); nil != e {
		return errors.New(e.Error() + ":>" + s)
	} else {
		return nil
	}
}

func (session *Session) Create(table *Table) error {
	s, args, e := session.getDialect().Create(table)
	if nil != e {
		return e
	}
	//log.Debugln("SQL:>>>", s)
	if _, e := session.db.Exec(s, args...); nil != e {
		return errors.New(e.Error() + ":>" + s)
	} else {
		return nil
	}
	//log.Debugln("SQL:>>>", s)
	//return e
}

//func (self *Session) Exec(s string,args...interface{}) (sql.Result, error) {
//    return self.Exec(s, args...)
//    //if _, e := self.db.Exec(s); nil != e {
//    //    return e
//    //} else {
//    //    return nil
//    //}
//}

func (session *Session) Close() {
	if session.tx != nil {
		_ = session.tx.Rollback()
	}
	//_ = session.db.Close()
}

func (session *Session) Table(table *Table, columns ...interface{}) QuerySet {
	qs := QuerySet{session: session, offset: -1, limit: -1}
	qs.table = table
	if len(columns) > 0 {
		for i, c := range columns {
			if qc, ok := c.(QueryColumn); ok {
				qs.queries = append(qs.queries, qc)
			} else if qcn, ok := c.(string); ok {
				if col, ok := qs.table.GetColumn(qcn); !ok {
					//panic("Invalid column name:" + qcn)
					qs.queries = append(qs.queries, QueryColumn{FieldName: qcn, Alias: fmt.Sprintf("aa%d", i)})
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

func (session *Session) Begin() error {
	if nil != session.tx {
		return errors.New("Already in Tx!")
	}
	tx, err := session.db.Begin()
	if nil == err {
		session.tx = tx
	}
	return err
}

func (session *Session) Commit() error {
	if nil == session.tx {
		return errors.New("not open Tx!")
	}
	err := session.tx.Commit()
	if nil == err {
		session.tx = nil
	}
	return err
}

func (session *Session) Rollback() error {
	if nil == session.tx {
		return errors.New("not open Tx!")
	}
	err := session.tx.Rollback()
	session.tx = nil
	return err
}

func log_timing(t1 time.Time, msg string, params ...interface{}) {
	t2 := time.Now()
	fmt.Println(t1, t2, t2.Sub(t1), msg, params)
}

func (session *Session) Exec(query string, args ...interface{}) (sql.Result, error) {
	if session.verbose {
		t1 := time.Now()
		defer log_timing(t1, "Session.Exec:", query)
	}
	if session.tx != nil {
		if session.verbose {
			//log.Debugln("Exec in Tx: ", query, args)
		}
		return session.tx.Exec(query, args...)
	} else {
		if session.verbose {
			//log.Debugln("Exec in DB: ", query, args)
		}
		return session.db.Exec(query, args...)
	}
	//return nil, nil
}

func (session *Session) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if session.verbose {
		t1 := time.Now()
		defer log_timing(t1, "Session.Query:", query)
	}
	if session.tx != nil {
		if session.verbose {
			//log.Debugln("Query in Tx: ", query, args)
		}
		return session.tx.Query(query, args...)
	} else {
		if session.verbose {
			//log.Debugln("Query in DB: ", query, args)
		}
		return session.db.Query(query, args...)
	}
	//return nil, nil
}

func (session *Session) QueryRow(query string, args ...interface{}) *sql.Row {
	if session.verbose {
		t1 := time.Now()
		defer log_timing(t1, "Session.doQueryRaw:", query)
	}
	if session.tx != nil {
		if session.verbose {
			//log.Debugln("QueryRow in Tx: ", query, args)
		}
		return session.tx.QueryRow(query, args...)
	} else {
		if session.verbose {
			//log.Debugln("QueryRow in Db: ", query, args)
		}
		return session.db.QueryRow(query, args...)
	}
	//return nil
}
