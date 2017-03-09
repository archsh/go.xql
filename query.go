package xql

import (
    "errors"
    log "github.com/Sirupsen/logrus"
    "reflect"
    "database/sql"
)



func (qc QueryColumn) String(as...bool) string {
    s := ""
    if qc.Function != "" {
        s = qc.Function+"("+qc.FieldName+")"
    }else{
        s = qc.FieldName
    }
    if qc.Alias != "" && len(as)>0 && as[0] {
        s = s + " AS " + qc.Alias
    }
    return s
}

type QuerySet struct {
    session *Session
    table   *Table
    queries []QueryColumn
    filters []QueryFilter
    orders  []QueryOrder
    offset  int64
    limit   int64
}


type XRow struct {
    row *sql.Row
    qs *QuerySet
}

func (self *XRow) Scan(dest ...interface{}) error {
    if nil == self.row {
        return errors.New("Nil row.")
    }
    if len(dest) < 1 {
        panic("Empty output!")
    }
    if len(dest) == 1 {
        d := dest[0]
        if reflect.TypeOf(d) == reflect.TypeOf(self.qs.table.Entity) {
            var outputs []interface{}
            r := reflect.ValueOf(d)
            for _, qc := range self.qs.queries {
                c, _ := self.qs.table.Columns[qc.FieldName]
                vp := r.Elem().FieldByName(c.PropertyName).Addr().Interface()
                outputs = append(outputs, vp)
            }
            return self.row.Scan(outputs...)
        }
    }
    return self.row.Scan(dest...)
}

type XRows struct {
    rows *sql.Rows
    qs *QuerySet
}

func (self *XRows) Scan(dest ...interface{}) error {
    if nil == self.rows {
        return errors.New("No rows.")
    }
    if len(dest) < 1 {
        panic("Empty output!")
    }
    if len(dest) == 1 {
        d := dest[0]
        if reflect.TypeOf(d) == reflect.TypeOf(self.qs.table.Entity) {
            var outputs []interface{}
            r := reflect.ValueOf(d)
            for _, qc := range self.qs.queries {
                c, _ := self.qs.table.Columns[qc.FieldName]
                vp := r.Elem().FieldByName(c.PropertyName).Addr().Interface()
                outputs = append(outputs, vp)
            }
            return self.rows.Scan(outputs...)
        }
    }
    return self.rows.Scan(dest...)
}

func (self *XRows) Next() bool {
    if nil == self.rows {
        return false
    }
    return self.rows.Next()
}

func (self *XRows) Close() {
    self.rows.Close()
    self.rows = nil
}


func makeQueryOrder(table *Table, s string) QueryOrder {
    qo := QueryOrder{}
    if s[:1] == "-" {
        qo.Type = ORDER_DESC
        qo.Field = s[1:]
    } else {
        qo.Type = ORDER_ASC
        qo.Field = s
    }
    return qo
}

func (self *QuerySet) Filter(cons ...interface{}) *QuerySet {
    for _, con := range cons {
        if vs, ok := con.(string); ok {
            self.filters = append(self.filters, QueryFilter{
                Field: vs,
            })
        }else if vm, ok := con.(map[string]interface{}); ok {
            for k, v := range vm {
                self.filters = append(self.filters, QueryFilter{
                    Field:k,
                    Value:v,
                    Operator: "=",
                })
            }
        }else if vf, ok := con.(QueryFilter); ok {
            self.filters = append(self.filters, vf)
        }
    }
    return self
}

func (self *QuerySet) OrderBy(orders ...interface{}) *QuerySet {
    for _, x := range orders {
        switch x.(type) {
        case string:
            qo := makeQueryOrder(self.table, x.(string))
            self.orders = append(self.orders, qo)
        case QueryOrder:
            self.orders = append(self.orders, x.(QueryOrder))
        default:
            panic("Not supported parameter type.")
        }
    }
    return self
}

func (self *QuerySet) Offset(offset int64) *QuerySet {
    self.offset = offset
    return self
}

func (self *QuerySet) Limit(limit int64) *QuerySet {
    self.limit = limit
    return self
}

func (self *QuerySet) Count(cols...string) (int64,error) {
    s, args, err := self.session.getStatementBuilder().Select(self.table,
        []QueryColumn{{Function:"COUNT", FieldName:"*"}},
        self.filters, nil, -1, -1)
    if nil != err {
        return 0, err
    }
    row := self.session.doQueryRow(s, args...)
    var n int64
    if e := row.Scan(&n); nil != e {
        return 0, e
    }
    return n, nil
}

func (self *QuerySet) All() (*XRows, error) {
    s, args, err := self.session.getStatementBuilder().Select(self.table, self.queries,
        self.filters, self.orders, self.offset, self.limit)
    if nil != err {
        return nil, err
    }
    rows, err := self.session.doQuery(s, args...)
    if nil != err {
        return nil, err
    }
    xrows := &XRows{rows:rows, qs:self}
    return xrows, nil
}


func (self *QuerySet) One() *XRow {
    s, args, err := self.session.getStatementBuilder().Select(self.table, self.queries,
        self.filters, self.orders, self.offset, 1)
    if nil != err {
        return nil
    }
    row := self.session.doQueryRow(s, args...)
    xrow := &XRow{row:row, qs:self}
    return xrow
}

func (self *QuerySet) Update(vals interface{}) (int64, error) {
    cols := []UpdateColumn{}
    if cm, ok := vals.(map[string]interface{}); ok {
        for k, v := range cm {
            uc := UpdateColumn{Field:k, Value:v, Operator:"="}
            cols = append(cols, uc)
        }
    }else if cx, ok := vals.([]UpdateColumn); ok {
        cols = cx
    }
    s, args, err := self.session.getStatementBuilder().Update(self.table, self.filters, cols...)
    if nil != err {
        return 0, err
    }
    var ret sql.Result
    ret, err = self.session.doExec(s, args...)
    if nil != err {
        return 0, err
    }else{
        rows, e := ret.RowsAffected()
        return rows, e
    }
    return 0, nil
}

func (self *QuerySet) Delete() (int64, error) {
    s, args, err := self.session.getStatementBuilder().Delete(self.table, self.filters)
    if nil != err {
        return 0, err
    }
    var ret sql.Result
    ret, err = self.session.doExec(s, args...)
    if nil != err {
        return 0, err
    }else{
        rows, e := ret.RowsAffected()
        return rows, e
    }
    return 0, nil
}

func (self *QuerySet) Insert(objs ...interface{}) (int64, error) {
    //log.Debugln("Insert:> ", objs)
    //log.Debugln("Table:> ", self.table)
    //var ret sql.Result
    var rows int64 = 0
    var cols []string
    if len(self.queries) > 0 {
        for _,x := range self.queries {
            cols = append(cols, x.FieldName)
        }
    }
    for _, obj := range objs {
        log.Debugln("reflect.TypeOf(obj)>", reflect.TypeOf(obj))
        log.Debugln("reflect.TypeOf(self.table.Entity)>", reflect.TypeOf(self.table.Entity))
        if reflect.TypeOf(obj) != reflect.TypeOf(self.table.Entity) {
           return 0, errors.New("Invalid data type.")
        }
        s, args, err := self.session.getStatementBuilder().Insert(self.table, obj, cols...)
        if nil != err {
            return 0, err
        }
        _, err = self.session.doExec(s, args...)
        if nil != err {
            return 0, err
        }else{
            //if nil != ret {
            //    n, e := ret.LastInsertId()
            //    if nil == e {
            //        rows += n
            //    }else{
            //        rows += 1
            //    }
            //    //log.Debugln("Insert:> ", n)
            //}else{
            //    rows += 1
            //}
            rows += 1
        }
    }
    return rows, nil
}
