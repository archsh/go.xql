package xql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

var pureFieldRex = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")

func isPureField(s string) bool {
	return pureFieldRex.MatchString(s)
}

func (qc QueryColumn) String(as ...bool) string {
	s := ""
	if qc.Function != "" {
		s = fmt.Sprintf(`%s("%s")`, qc.Function, qc.FieldName) //qc.Function+"("+qc.FieldName+")"
	} else if isPureField(qc.FieldName) {
		s = fmt.Sprintf(`"%s"`, qc.FieldName)
	} else {
		s = qc.FieldName
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
	lockFor string
	offset  int64
	limit   int64
}

type XRow struct {
	row *sql.Row
	qs  *QuerySet
}

func (xr *XRow) Scan(dest ...interface{}) error {
	if nil == xr.row {
		return errors.New("nil row")
	}
	if len(dest) < 1 {
		panic("Empty output!")
	}
	if len(dest) == 1 {
		d := dest[0]
		if reflect.TypeOf(d) == reflect.TypeOf(xr.qs.table.entity) {
			var outputs []interface{}
			r := reflect.ValueOf(d)
			for _, qc := range xr.qs.queries {
				c, _ := xr.qs.table.GetColumn(qc.FieldName)
				//fmt.Printf("> Scan to field '%s' to '%s' .\n", qc.FieldName, c.ElemName)
				vp := r.Elem().FieldByName(c.ElemName).Addr().Interface()
				outputs = append(outputs, vp)
			}
			return xr.row.Scan(outputs...)
		}
	}
	return xr.row.Scan(dest...)
}

type XRows struct {
	rows *sql.Rows
	qs   *QuerySet
}

func (xr *XRows) Scan(dest ...interface{}) error {
	if nil == xr.rows {
		return errors.New("no rows")
	}
	if len(dest) < 1 {
		panic("Empty output!")
	}
	if len(dest) == 1 {
		d := dest[0]
		if reflect.TypeOf(d) == reflect.TypeOf(xr.qs.table.entity) {
			var outputs []interface{}
			r := reflect.ValueOf(d)
			for _, qc := range xr.qs.queries {
				c, _ := xr.qs.table.GetColumn(qc.FieldName)
				vp := r.Elem().FieldByName(c.ElemName).Addr().Interface()
				outputs = append(outputs, vp)
			}
			return xr.rows.Scan(outputs...)
		}
	}
	return xr.rows.Scan(dest...)
}

func (xr *XRows) Next() bool {
	if nil == xr.rows {
		return false
	}
	return xr.rows.Next()
}

func (xr *XRows) Close() {
	if xr.rows != nil {
		xr.rows.Close()
		xr.rows = nil
	}
}

func makeQueryOrder(table *Table, s string) QueryOrder {
	qo := QueryOrder{}
	if s[:1] == "-" {
		qo.Type = OrderDesc
		qo.Field = s[1:]
	} else {
		qo.Type = OrderAsc
		qo.Field = s
	}
	return qo
}

func Where(field string, val interface{}, ops ...string) QueryFilter {
	f := QueryFilter{Field: field, Value: val, Operator: "="}
	if len(ops) > 0 {
		f.Operator = ops[0]
		if len(ops) > 1 {
			f.Function = ops[1]
		}
	}
	return f
}

func (qs QuerySet) Where(field string, val interface{}, ops ...string) QuerySet {
	f := QueryFilter{Field: field, Value: val, Operator: "="}
	if len(ops) > 0 {
		f.Operator = ops[0]
		if len(ops) > 1 {
			f.Function = ops[1]
		}
	}
	qs.filters = append(qs.filters, f)
	return qs
}

func (qs QuerySet) And(field string, val interface{}, ops ...string) QuerySet {
	f := QueryFilter{Field: field, Value: val, Operator: "="}
	if len(ops) > 0 {
		f.Operator = ops[0]
		if len(ops) > 1 {
			f.Function = ops[1]
		}
	}
	qs.filters = append(qs.filters, f)
	return qs
}

func (qs QuerySet) Or(field string, val interface{}, ops ...string) QuerySet {
	f := QueryFilter{Field: field, Value: val, Operator: "=", Condition: ConditionOr}
	if len(ops) > 0 {
		f.Operator = ops[0]
		if len(ops) > 1 {
			f.Function = ops[1]
		}
	}
	qs.filters = append(qs.filters, f)
	return qs
}

func (qs QuerySet) LockFor(s string) QuerySet {
	qs.lockFor = s
	return qs
}

func (qs QuerySet) Filter(cons ...interface{}) QuerySet {
	for _, con := range cons {
		if vs, ok := con.(string); ok {
			qs.filters = append(qs.filters, QueryFilter{
				Field: vs,
			})
		} else if vm, ok := con.(map[string]interface{}); ok {
			for k, v := range vm {
				qs.filters = append(qs.filters, QueryFilter{
					Field:    k,
					Value:    v,
					Operator: "=",
				})
			}
		} else if vf, ok := con.(*QueryFilter); ok {
			qs.filters = append(qs.filters, *vf)
		} else if vf, ok := con.(QueryFilter); ok {
			qs.filters = append(qs.filters, vf)
		} else {
			panic("Unknow Filter!")
		}
	}
	return qs
}

func (qs QuerySet) OrderBy(orders ...interface{}) QuerySet {
	for _, x := range orders {
		switch x.(type) {
		case string:
			qo := makeQueryOrder(qs.table, x.(string))
			qs.orders = append(qs.orders, qo)
		case *QueryOrder:
			qs.orders = append(qs.orders, *x.(*QueryOrder))
		case QueryOrder:
			qs.orders = append(qs.orders, x.(QueryOrder))
		default:
			panic("Not supported parameter type.")
		}
	}
	return qs
}

func (qs QuerySet) Offset(offset int64) QuerySet {
	qs.offset = offset
	return qs
}

func (qs QuerySet) Limit(limit int64) QuerySet {
	qs.limit = limit
	return qs
}

func (qs QuerySet) Count(cols ...string) (int64, error) {
	var fieldName string
	if len(cols) > 0 {
		fieldName = cols[0]
	} else if len(qs.table.primaryKeys) > 0 {
		fieldName = qs.table.primaryKeys[0].FieldName
	} else {
		fieldName = qs.table.columns[0].FieldName
	}
	s, args, err := qs.session.getDialect().Select(qs.table,
		[]QueryColumn{{Function: "COUNT", FieldName: fieldName}},
		qs.filters, nil, qs.lockFor, -1, -1)
	if nil != err {
		return 0, err
	}
	row := qs.session.QueryRow(s, args...)
	var n int64
	if e := row.Scan(&n); nil != e {
		return 0, e
	}
	return n, nil
}

func (qs QuerySet) All() (*XRows, error) {
	if len(qs.queries) < 1 {
		for _, col := range qs.table.mColumns {
			qs.queries = append(qs.queries, QueryColumn{FieldName: col.FieldName, Alias: col.FieldName})
		}
	}
	s, args, err := qs.session.getDialect().Select(qs.table, qs.queries,
		qs.filters, qs.orders, qs.lockFor, qs.offset, qs.limit)
	if nil != err {
		return nil, err
	}
	rows, err := qs.session.Query(s, args...)
	if nil != err {
		return nil, err
	}
	xrows := &XRows{rows: rows, qs: &qs}
	return xrows, nil
}

func (qs QuerySet) One() *XRow {
	if len(qs.queries) < 1 {
		for _, col := range qs.table.mColumns {
			qs.queries = append(qs.queries, QueryColumn{FieldName: col.FieldName, Alias: col.FieldName})
		}
	}
	s, args, err := qs.session.getDialect().Select(qs.table, qs.queries,
		qs.filters, qs.orders, qs.lockFor, qs.offset, 1)
	if nil != err {
		return nil
	}
	//fmt.Println("One:>", s, args)
	row := qs.session.QueryRow(s, args...)
	xrow := &XRow{row: row, qs: &qs}
	return xrow
}

func (qs QuerySet) Get(pks ...interface{}) *XRow {
	if len(qs.queries) < 1 {
		for _, col := range qs.table.mColumns {
			qs.queries = append(qs.queries, QueryColumn{FieldName: col.FieldName, Alias: col.FieldName})
		}
	}
	if len(pks) != len(qs.table.primaryKeys) {
		panic("Primary Key number not match!")
	}
	qs.filters = make([]QueryFilter, 0, len(pks))
	for i, pk := range pks {
		filter := QueryFilter{
			Field:    qs.table.primaryKeys[i].FieldName,
			Value:    pk,
			Operator: "=",
		}
		qs.filters = append(qs.filters, filter)
	}
	s, args, err := qs.session.getDialect().Select(qs.table, qs.queries,
		qs.filters, qs.orders, qs.lockFor, qs.offset, 1)
	if nil != err {
		return nil
	}
	row := qs.session.QueryRow(s, args...)
	xrow := &XRow{row: row, qs: &qs}
	return xrow
}

func (qs QuerySet) Update(vals interface{}) (int64, error) {
	var cols []UpdateColumn
	//fmt.Println("Update:>", qs.table.mColumns)
	if cm, ok := vals.(map[string]interface{}); ok {
		for k, v := range cm {
			var fk string
			if c, ok := qs.table.xColumns[k]; ok {
				fk = c.FieldName
			} else if c, ok := qs.table.mColumns[k]; ok {
				fk = c.FieldName
			} else if c, ok := qs.table.jColumns[k]; ok {
				fk = c.FieldName
			} else {
				return 0, errors.New("Invalid column:" + k)
			}
			uc := UpdateColumn{Field: fk, Value: v, Operator: "="}
			cols = append(cols, uc)
		}
	} else if cx, ok := vals.([]UpdateColumn); ok {
		cols = cx
	} else if reflect.TypeOf(vals) == reflect.TypeOf(qs.table.entity) {
		r := reflect.ValueOf(vals)
		for _, col := range qs.table.columns {
			if col.PrimaryKey {
				continue
			}
			fv := reflect.Indirect(r).FieldByName(col.ElemName)
			if (fv.Kind() == reflect.Ptr && fv.IsNil()) || reflect.Zero(fv.Type()) == fv {
				continue
			} else {
				cols = append(cols, UpdateColumn{Field: col.FieldName, Operator: "=", Value: fv.Interface()})
			}
		}
	}
	s, args, err := qs.session.getDialect().Update(qs.table, qs.filters, cols...)
	if nil != err {
		return 0, err
	}
	//fmt.Println(">>>Update:", s, args)
	var ret sql.Result
	ret, err = qs.session.Exec(s, args...)
	if nil != err {
		return 0, err
	} else {
		rows, e := ret.RowsAffected()
		return rows, e
	}
	return 0, nil
}

func (qs QuerySet) Delete() (int64, error) {
	s, args, err := qs.session.getDialect().Delete(qs.table, qs.filters)
	if nil != err {
		return 0, err
	}
	var ret sql.Result
	ret, err = qs.session.Exec(s, args...)
	if nil != err {
		return 0, err
	} else {
		rows, e := ret.RowsAffected()
		return rows, e
	}
	return 0, nil
}

func (qs QuerySet) InsertWithInsertedId(obj interface{}, idname string, id interface{}) error {
	var cols []string
	if len(qs.queries) > 0 {
		for _, x := range qs.queries {
			cols = append(cols, x.FieldName)
		}
	}
	if reflect.TypeOf(obj) != reflect.TypeOf(qs.table.entity) {
		return errors.New(fmt.Sprintf("Invalid data type: %s <> %s", reflect.TypeOf(obj).String(),
			reflect.TypeOf(qs.table.entity).String()))
	}
	if pobj, ok := obj.(TablePreInsert); ok {
		pobj.PreInsert(qs.table, qs.session)
	}
	s, args, err := qs.session.getDialect().InsertWithInsertedId(qs.table, obj, idname, cols...)
	if nil != err {
		return err
	}
	//fmt.Println("Insert SQL:>", s, args)
	err = qs.session.QueryRow(s, args...).Scan(id)
	if nil != err {
		//fmt.Println(">>>Insert SQL:>", s, args, err)
		return err
	} else {
		if pobj, ok := obj.(TablePostInsert); ok {
			pobj.PostInsert(qs.table, qs.session)
		}
	}
	return nil
}

func (qs QuerySet) Insert(objs ...interface{}) (int64, error) {
	var rows int64 = 0
	var cols []string
	if len(qs.queries) > 0 {
		for _, x := range qs.queries {
			cols = append(cols, x.FieldName)
		}
	}
	for _, obj := range objs {
		if reflect.TypeOf(obj) != reflect.TypeOf(qs.table.entity) {
			return 0, errors.New(fmt.Sprintf("Invalid data type: %s <> %s", reflect.TypeOf(obj).String(),
				reflect.TypeOf(qs.table.entity).String()))
		}
		if pobj, ok := obj.(TablePreInsert); ok {
			pobj.PreInsert(qs.table, qs.session)
		}
		s, args, err := qs.session.getDialect().Insert(qs.table, obj, cols...)
		if nil != err {
			return 0, err
		}
		//fmt.Println("Insert SQL:>", s, args)
		_, err = qs.session.Exec(s, args...)
		if nil != err {
			//fmt.Println(">>>Insert SQL:>", s, args, err)
			return 0, err
		} else {
			if pobj, ok := obj.(TablePostInsert); ok {
				pobj.PostInsert(qs.table, qs.session)
			}
			rows += 1
		}
	}
	return rows, nil
}
