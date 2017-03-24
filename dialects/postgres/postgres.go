package postgres
import (
    "github.com/archsh/go.xql"
    "fmt"
    "reflect"
    "strings"
)

type PostgresDialect struct {


}

func (pb PostgresDialect) Create(t *xql.Table, options ...interface{}) (s string, args []interface{}, err error) {
    return
}


func (pb PostgresDialect) Select(t *xql.Table, cols []xql.QueryColumn, filters []xql.QueryFilter, orders []xql.QueryOrder, offset int64, limit int64)  (s string, args []interface{}, err error) {
    var colnames []string
    for _,x := range cols {
        colnames = append(colnames, x.String())
    }
    s = fmt.Sprintf("SELECT %s FROM ", strings.Join(colnames,","))
    if t.Schema != "" {
        s += t.Schema+"."
    }
    s += t.TableName
    var n int
    for i, f := range filters {
        var cause string
        switch f.Condition {
        case xql.CONDITION_AND:
            cause = "AND"
        case xql.CONDITION_OR:
            cause = "OR"
        }
        if i == 0 {
            cause = "WHERE"
        }
        if f.Operator == "" {
            s = fmt.Sprintf(`%s %s %s`, s, cause, f.Field)
        } else if f.Reversed {
            n += 1
            if f.Function != "" {
                s = fmt.Sprintf(`%s %s %s($%d) %s %s`, s, cause, f.Function, n, f.Operator, f.Field)
            }else{
                s = fmt.Sprintf(`%s %s $%d %s %s`, s, cause, n, f.Operator, f.Field)
            }

            args = append(args, f.Value)
        } else {
            n += 1
            if f.Function != "" {
                s = fmt.Sprintf(`%s %s %s %s %s($%d)`, s, cause, f.Field, f.Operator, f.Function, n)
            }else{
                s = fmt.Sprintf(`%s %s %s %s $%d`, s, cause, f.Field, f.Operator, n)
            }

            args = append(args, f.Value)
        }
    }
    var s_orders []string
    for _, o := range orders {
        switch o.Type {
        case xql.ORDER_ASC:
            s_orders = append(s_orders, fmt.Sprintf(`%s ASC`, o.Field))
        case xql.ORDER_DESC:
            s_orders = append(s_orders, fmt.Sprintf(`%s DESC`, o.Field))
        }
    }
    if len(s_orders) > 0 {
        s = fmt.Sprintf(`%s ORDER BY %s`, s, strings.Join(s_orders, ","))
    }
    if offset >= 0 {
        s = fmt.Sprintf(`%s OFFSET %d`, s, offset)
    }

    if limit >= 0 {
        s = fmt.Sprintf(`%s LIMIT %d`, s, limit)
    }
    return
}

func (pb PostgresDialect) Insert(t *xql.Table, obj interface{}, col...string) (s string, args []interface{}, err error) {
    s = "INSERT INTO "
    if t.Schema != "" {
        s += t.Schema+"."
    }
    s += t.TableName
    var cols []string
    var vals []string
    i := 0
    r := reflect.ValueOf(obj)
    if len(col) > 0 {
        for _, n := range col {
            v, ok := t.Columns[n]
            if ! ok {
                continue
            }
            //fmt.Println("POSTGRES Insert>1>>>",n,v.Auto,v.PrimaryKey,v)
            i += 1
            cols = append(cols, n)
            vals = append(vals, fmt.Sprintf("$%d", i))
            fv := reflect.Indirect(r).FieldByName(v.PropertyName).Interface()
            args = append(args, fv)
        }
    }else{
        for k, v := range t.Columns {
            //fmt.Println("POSTGRES Insert>2>>>",k,v.Auto,v.PrimaryKey,v)
            if v.Auto {
                continue
            }
            i += 1
            cols = append(cols, k)
            vals = append(vals, fmt.Sprintf("$%d", i))
            fv := reflect.Indirect(r).FieldByName(v.PropertyName).Interface()
            args = append(args, fv)
        }
    }

    s = fmt.Sprintf("%s (%s) VALUES(%s)", s, strings.Join(cols,","), strings.Join(vals,","))
    return
}

func (pb PostgresDialect) Update(t *xql.Table, filters []xql.QueryFilter, cols ...xql.UpdateColumn) (s string, args []interface{}, err error) {
    s = "UPDATE "
    if t.Schema != "" {
        s += t.Schema+"."
    }
    s += t.TableName
    if len(cols) < 1 {
        panic("Empty Update Columns!!!")
    }
    var n int
    for i, uc := range cols {
        n += 1
        if i == 0 {
            //s = s + " SET "
            s = fmt.Sprintf(`%s SET %s=$%d`,s, uc.Field, n)
        }else{
            s = fmt.Sprintf(`%s, %s=$%d`,s, uc.Field, n)
        }


        args = append(args, uc.Value)
    }

    for i, f := range filters {
        var cause string
        switch f.Condition {
        case xql.CONDITION_AND:
            cause = "AND"
        case xql.CONDITION_OR:
            cause = "OR"
        }
        if i == 0 {
            cause = "WHERE"
        }
        if f.Operator == "" {
            s = fmt.Sprintf("%s %s %s", s, cause, f.Field)
        } else if f.Reversed {
            n += 1
            s = fmt.Sprintf("%s %s $%d %s %s", s, cause, n, f.Operator, f.Field)
            args = append(args, f.Value)
        } else {
            n += 1
            s = fmt.Sprintf("%s %s %s %s $%d", s, cause, f.Field, f.Operator, n)
            args = append(args, f.Value)
        }
    }
    return
}

func (pb PostgresDialect) Delete(t *xql.Table, filters []xql.QueryFilter) (s string, args []interface{}, err error) {
    s = "DELETE FROM "
    if t.Schema != "" {
        s += t.Schema+"."
    }
    s += t.TableName
    var n int
    for i, f := range filters {
        var cause string
        switch f.Condition {
        case xql.CONDITION_AND:
            cause = "AND"
        case xql.CONDITION_OR:
            cause = "OR"
        }
        if i == 0 {
            cause = "WHERE"
        }
        if f.Operator == "" {
            s = fmt.Sprintf("%s %s %s", s, cause, f.Field)
        } else if f.Reversed {
            n += 1
            s = fmt.Sprintf("%s %s $%d %s %s", s, cause, n, f.Operator, f.Field)
            args = append(args, f.Value)
        } else {
            n += 1
            s = fmt.Sprintf("%s %s %s %s $%d", s, cause, f.Field, f.Operator, n)
            args = append(args, f.Value)
        }
    }
    return
}

func init() {
    xql.RegisterDialect("postgres", &PostgresDialect{})
}