package xql

import (
    "errors"
    "reflect"
    "database/sql"
    "fmt"
)

func (qc QueryColumn) String(as ...bool) string {
    s := ""
    if qc.Function != "" {
        s = fmt.Sprintf(`%s("%s")`, qc.Function, qc.FieldName) //qc.Function+"("+qc.FieldName+")"
    } else {
        s = fmt.Sprintf(`"%s"`, qc.FieldName)
    }
    if qc.Alias != "" && len(as) > 0 && as[0] {
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
    qs  *QuerySet
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
        if reflect.TypeOf(d) == reflect.TypeOf(self.qs.table.entity) {
            var outputs []interface{}
            r := reflect.ValueOf(d)
            for _, qc := range self.qs.queries {
                c, _ := self.qs.table.GetColumn(qc.FieldName)
                //fmt.Printf("> Scan to field '%s' to '%s' .\n", qc.FieldName, c.ElemName)
                vp := r.Elem().FieldByName(c.ElemName).Addr().Interface()
                outputs = append(outputs, vp)
            }
            return self.row.Scan(outputs...)
        }
    }
    return self.row.Scan(dest...)
}

type XRows struct {
    rows *sql.Rows
    qs   *QuerySet
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
        if reflect.TypeOf(d) == reflect.TypeOf(self.qs.table.entity) {
            var outputs []interface{}
            r := reflect.ValueOf(d)
            for _, qc := range self.qs.queries {
                c, _ := self.qs.table.GetColumn(qc.FieldName)
                vp := r.Elem().FieldByName(c.ElemName).Addr().Interface()
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
    if self.rows != nil {
        self.rows.Close()
        self.rows = nil
    }
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

func (self QuerySet) Where(field string, val interface{}, ops ...string) QuerySet {
    f := QueryFilter{Field: field, Value: val, Operator: "="}
    if len(ops) > 0 {
        f.Operator = ops[0]
        if len(ops) > 1 {
            f.Function = ops[1]
        }
    }
    self.filters = append(self.filters, f)
    return self
}

func (self QuerySet) And(field string, val interface{}, ops ...string) QuerySet {
    f := QueryFilter{Field: field, Value: val, Operator: "="}
    if len(ops) > 0 {
        f.Operator = ops[0]
        if len(ops) > 1 {
            f.Function = ops[1]
        }
    }
    self.filters = append(self.filters, f)
    return self
}

func (self QuerySet) Or(field string, val interface{}, ops ...string) QuerySet {
    f := QueryFilter{Field: field, Value: val, Operator: "=", Condition: CONDITION_OR}
    if len(ops) > 0 {
        f.Operator = ops[0]
        if len(ops) > 1 {
            f.Function = ops[1]
        }
    }
    self.filters = append(self.filters, f)
    return self
}

func (self QuerySet) Filter(cons ...interface{}) QuerySet {
    for _, con := range cons {
        if vs, ok := con.(string); ok {
            self.filters = append(self.filters, QueryFilter{
                Field: vs,
            })
        } else if vm, ok := con.(map[string]interface{}); ok {
            for k, v := range vm {
                self.filters = append(self.filters, QueryFilter{
                    Field:    k,
                    Value:    v,
                    Operator: "=",
                })
            }
        } else if vf, ok := con.(*QueryFilter); ok {
            self.filters = append(self.filters, *vf)
        } else if vf, ok := con.(QueryFilter); ok {
            self.filters = append(self.filters, vf)
        } else {
            panic("Unknow Filter!")
        }
    }
    return self
}

func (self QuerySet) OrderBy(orders ...interface{}) QuerySet {
    for _, x := range orders {
        switch x.(type) {
        case string:
            qo := makeQueryOrder(self.table, x.(string))
            self.orders = append(self.orders, qo)
        case *QueryOrder:
            self.orders = append(self.orders, *x.(*QueryOrder))
        case QueryOrder:
            self.orders = append(self.orders, x.(QueryOrder))
        default:
            panic("Not supported parameter type.")
        }
    }
    return self
}

func (self QuerySet) Offset(offset int64) QuerySet {
    self.offset = offset
    return self
}

func (self QuerySet) Limit(limit int64) QuerySet {
    self.limit = limit
    return self
}

func (self QuerySet) Count(cols ...string) (int64, error) {
    var fieldname string
    if len(cols) > 0 {
        fieldname = cols[0]
    } else if len(self.table.primary_keys) > 0 {
        fieldname = self.table.primary_keys[0].FieldName
    } else {
        fieldname = self.table.columns[0].FieldName
    }
    s, args, err := self.session.getDialect().Select(self.table,
        []QueryColumn{{Function: "COUNT", FieldName: fieldname}},
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

func (self QuerySet) All() (*XRows, error) {
    if len(self.queries) < 1 {
        for _, col := range self.table.m_columns {
            self.queries = append(self.queries, QueryColumn{FieldName: col.FieldName, Alias: col.FieldName})
        }
    }
    s, args, err := self.session.getDialect().Select(self.table, self.queries,
        self.filters, self.orders, self.offset, self.limit)
    if nil != err {
        return nil, err
    }
    rows, err := self.session.doQuery(s, args...)
    if nil != err {
        return nil, err
    }
    xrows := &XRows{rows: rows, qs: &self}
    return xrows, nil
}

func (self QuerySet) One() *XRow {
    if len(self.queries) < 1 {
        for _, col := range self.table.m_columns {
            self.queries = append(self.queries, QueryColumn{FieldName: col.FieldName, Alias: col.FieldName})
        }
    }
    s, args, err := self.session.getDialect().Select(self.table, self.queries,
        self.filters, self.orders, self.offset, 1)
    if nil != err {
        return nil
    }
    //fmt.Println("One:>", s, args)
    row := self.session.doQueryRow(s, args...)
    xrow := &XRow{row: row, qs: &self}
    return xrow
}

func (self QuerySet) Get(pks ...interface{}) *XRow {
    if len(self.queries) < 1 {
        for _, col := range self.table.m_columns {
            self.queries = append(self.queries, QueryColumn{FieldName: col.FieldName, Alias: col.FieldName})
        }
    }
    if len(pks) != len(self.table.primary_keys) {
        panic("Primary Key number not match!")
    }
    self.filters = make([]QueryFilter,0,len(pks))
    for i, pk := range pks {
        filter := QueryFilter{
            Field: self.table.primary_keys[i].FieldName,
            Value: pk,
            Operator: "=",
        }
        self.filters = append(self.filters, filter)
    }
    s, args, err := self.session.getDialect().Select(self.table, self.queries,
        self.filters, self.orders, self.offset, 1)
    if nil != err {
        return nil
    }
    row := self.session.doQueryRow(s, args...)
    xrow := &XRow{row: row, qs: &self}
    return xrow
}

func (self QuerySet) Update(vals interface{}) (int64, error) {
    cols := []UpdateColumn{}
    //fmt.Println("Update:>", self.table.m_columns)
    if cm, ok := vals.(map[string]interface{}); ok {
        for k, v := range cm {
            var fk string
            if c, ok := self.table.x_columns[k]; ok {
                fk = c.FieldName
            } else if c, ok := self.table.m_columns[k]; ok {
                fk = c.FieldName
            } else if c, ok := self.table.j_columns[k]; ok {
                fk = c.FieldName
            } else {
                return 0, errors.New("Invalid column:" + k)
            }
            uc := UpdateColumn{Field: fk, Value: v, Operator: "="}
            cols = append(cols, uc)
        }
    } else if cx, ok := vals.([]UpdateColumn); ok {
        cols = cx
    } else if reflect.TypeOf(vals) == reflect.TypeOf(self.table.entity) {
        r := reflect.ValueOf(vals)
        for _, col := range self.table.columns {
            if col.PrimaryKey {
                continue
            }
            fv := reflect.Indirect(r).FieldByName(col.ElemName)
            if ( fv.Kind() == reflect.Ptr && fv.IsNil() ) || reflect.Zero(fv.Type()) == fv {
                continue
            }else{
                cols = append(cols, UpdateColumn{Field:col.FieldName, Operator:"=", Value:fv.Interface()})
            }
        }
    }
    s, args, err := self.session.getDialect().Update(self.table, self.filters, cols...)
    if nil != err {
        return 0, err
    }
    //fmt.Println(">>>Update:", s, args)
    var ret sql.Result
    ret, err = self.session.doExec(s, args...)
    if nil != err {
        return 0, err
    } else {
        rows, e := ret.RowsAffected()
        return rows, e
    }
    return 0, nil
}

func (self QuerySet) Delete() (int64, error) {
    s, args, err := self.session.getDialect().Delete(self.table, self.filters)
    if nil != err {
        return 0, err
    }
    var ret sql.Result
    ret, err = self.session.doExec(s, args...)
    if nil != err {
        return 0, err
    } else {
        rows, e := ret.RowsAffected()
        return rows, e
    }
    return 0, nil
}


func (self QuerySet) InsertWithInsertedId(obj interface{}, idname string, id interface{}) error {
    var cols []string
    if len(self.queries) > 0 {
        for _, x := range self.queries {
            cols = append(cols, x.FieldName)
        }
    }
    if reflect.TypeOf(obj) != reflect.TypeOf(self.table.entity) {
        return errors.New(fmt.Sprintf("Invalid data type: %s <> %s",reflect.TypeOf(obj).String(),
            reflect.TypeOf(self.table.entity).String()))
    }
    if pobj, ok := obj.(TablePreInsert); ok {
        pobj.PreInsert(self.table, self.session)
    }
    s, args, err := self.session.getDialect().InsertWithInsertedId(self.table, obj, idname, cols...)
    if nil != err {
        return err
    }
    //fmt.Println("Insert SQL:>", s, args)
    err = self.session.doQueryRow(s, args...).Scan(id)
    if nil != err {
        //fmt.Println(">>>Insert SQL:>", s, args, err)
        return err
    } else {
        if pobj, ok := obj.(TablePostInsert); ok {
            pobj.PostInsert(self.table, self.session)
        }
    }
    return nil
}

func (self QuerySet) Insert(objs ...interface{}) (int64, error) {
    var rows int64 = 0
    var cols []string
    if len(self.queries) > 0 {
        for _, x := range self.queries {
            cols = append(cols, x.FieldName)
        }
    }
    for _, obj := range objs {
        if reflect.TypeOf(obj) != reflect.TypeOf(self.table.entity) {
            return 0, errors.New(fmt.Sprintf("Invalid data type: %s <> %s",reflect.TypeOf(obj).String(),
                reflect.TypeOf(self.table.entity).String()))
        }
        if pobj, ok := obj.(TablePreInsert); ok {
            pobj.PreInsert(self.table, self.session)
        }
        s, args, err := self.session.getDialect().Insert(self.table, obj, cols...)
        if nil != err {
            return 0, err
        }
        //fmt.Println("Insert SQL:>", s, args)
        _, err = self.session.doExec(s, args...)
        if nil != err {
            //fmt.Println(">>>Insert SQL:>", s, args, err)
            return 0, err
        } else {
            if pobj, ok := obj.(TablePostInsert); ok {
                pobj.PostInsert(self.table, self.session)
            }
            rows += 1
        }
    }
    return rows, nil
}
