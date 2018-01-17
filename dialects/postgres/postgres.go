package postgres

import (
    "github.com/archsh/go.xql"
    "fmt"
    "reflect"
    "strings"
    "errors"
    "database/sql"
)

type PostgresDialect struct {
}

/*
-- Table: metas_vod_albums

-- DROP TABLE metas_vod_albums;

CREATE TABLE metas_vod_albums
(
  id uuid NOT NULL,
  album_name character varying(256) NOT NULL,
  idx integer NOT NULL,
  is_series boolean NOT NULL,
  film_total integer NOT NULL,
  film_count integer NOT NULL,
  remark character varying(256) NOT NULL,
  issue_date date NOT NULL,
  publish_time timestamp without time zone NOT NULL,
  score numeric(3,1) NOT NULL,
  content_rating character varying(16),
  album_pic character varying(256) NOT NULL,
  imdb_id character varying(24),
  cp_ref_id character varying(64),
  description text NOT NULL,
  status smallint NOT NULL,
  created timestamp without time zone NOT NULL,
  updated timestamp without time zone NOT NULL,
  provider_id uuid NOT NULL,
  alias character varying(256) NOT NULL,
  CONSTRAINT metas_vod_albums_pkey PRIMARY KEY (id),
  CONSTRAINT metas_vod_albums_provider_id_fkey FOREIGN KEY (provider_id)
      REFERENCES metas_providers (id) MATCH SIMPLE
      ON UPDATE CASCADE ON DELETE CASCADE
)
WITH (
  OIDS=FALSE
);
ALTER TABLE metas_vod_albums
  OWNER TO postgres;

-- Index: ix_public_metas_vod_albums_album_name

-- DROP INDEX ix_public_metas_vod_albums_album_name;

CREATE INDEX ix_public_metas_vod_albums_album_name
  ON metas_vod_albums
  USING btree
  (album_name COLLATE pg_catalog."default");

-- Index: ix_public_metas_vod_albums_idx

-- DROP INDEX ix_public_metas_vod_albums_idx;

CREATE INDEX ix_public_metas_vod_albums_idx
  ON metas_vod_albums
  USING btree
  (idx);
 */

 func makeInlineConstraint(c... *xql.Constraint) string {
     constraints := []string{}
     for _, x := range c {
         switch x.Type {
         case xql.CONSTRAINT_NOT_NULL:
             constraints = append(constraints, "NOT NULL")
         case xql.CONSTRAINT_UNIQUE:
             constraints = append(constraints, "UNIQUE")
         case xql.CONSTRAINT_CHECK:
             constraints = append(constraints, fmt.Sprintf("CHECK (%s)", x.Statement))
         //case xql.CONSTRAINT_EXCLUDE:
             //constraints = append(constraints, "NOT NUL")
         case xql.CONSTRAINT_FOREIGNKEY:
             ondelete := ""
             onupdate := ""
             if x.OnDelete != "" {
                 ondelete = "ON DELETE "+x.OnDelete
             }
             if x.OnUpdate != "" {
                 onupdate = "ON UPDATE "+x.OnUpdate
             }
             xs := strings.Split(x.Statement, ".")
             if len(xs) > 1 {
                 tt := strings.Join(xs[:len(xs)-1],".")
                 tc := xs[len(xs)-1]
                 constraints = append(constraints, fmt.Sprintf("REFERENCES %s (%s) %s %s",tt, tc, onupdate, ondelete))
             }else{
                 constraints = append(constraints, fmt.Sprintf("REFERENCES %s %s %s",xs[0], onupdate, ondelete))
             }

         case xql.CONSTRAINT_PRIMARYKEY:
             constraints = append(constraints, "PRIMARY KEY")
         }
     }
     if len(constraints) < 1 {
         return ""
     }
     return strings.Join(constraints, " ")
 }

 func makeConstraints(t *xql.Table, idx int, c... *xql.Constraint) (ret []string) {
     for _, x := range c {
         fields := []string{}
         for _,cc := range x.Columns {
             fields = append(fields, cc.FieldName)
         }
         field_str := strings.Join(fields, ",")
         name_str := fmt.Sprintf("%s_%s", t.BaseTableName(), strings.Join(fields,"_"))
         switch x.Type {
         //case xql.CONSTRAINT_NOT_NULL:
         //    ret = append(ret, fmt.Sprintf("NOT NUL"))
         case xql.CONSTRAINT_UNIQUE:
             ret = append(ret, fmt.Sprintf("CONSTRAINT %s_unique UNIQUE (%s)", name_str, field_str))
         case xql.CONSTRAINT_CHECK:
             ret = append(ret, fmt.Sprintf("CONSTRAINT %s_check CHECK (%s)", name_str, x.Statement))
         case xql.CONSTRAINT_EXCLUDE:
             ret = append(ret, fmt.Sprintf("CONSTRAINT %s_exclude EXCLUDE USING %s", name_str, x.Statement))
         case xql.CONSTRAINT_FOREIGNKEY:
             ondelete := ""
             onupdate := ""
             if x.OnDelete != "" {
                 ondelete = "ON DELETE "+x.OnDelete
             }
             if x.OnUpdate != "" {
                 onupdate = "ON UPDATE "+x.OnUpdate
             }
             ret = append(ret,
                 fmt.Sprintf("CONSTRAINT %s_fkey FOREIGN KEY (%s) REFERENCES %s %s %s",
                     name_str, field_str, x.Statement, onupdate, ondelete))
         case xql.CONSTRAINT_PRIMARYKEY:
             ret = append(ret, fmt.Sprintf("CONSTRAINT %s_pkey PRIMARY KEY (%s)", name_str, field_str))
         }
     }
     return
 }

 func makeIndexes(t *xql.Table, idx int, i... *xql.Index) (ret []string) {
     for _, ii := range i {
         fs := []string{}
         for _, c := range ii.Columns {
             fs = append(fs, c.FieldName)
         }
         tp := ""
         switch ii.Type {
         case xql.INDEX_B_TREE:
             tp = "USING btree"
         case xql.INDEX_HASH:
             tp = "USING hash"
         case xql.INDEX_BRIN:
             tp = "USING brin"
         case xql.INDEX_GIST:
             tp = "USING gist"
         case xql.INDEX_SP_GIST:
             tp = "USING sp-gist"
         case xql.INDEX_GIN:
             tp = "USING gin"
         }
         // CREATE INDEX test2_mm_idx ON test2 (major, minor);
         s := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s %s (\"%s\");",ii.Name , t.BaseTableName(), tp, strings.Join(fs, ","))
         ret = append(ret, s)
     }
     return
 }

 // Drop
 // Implement the IDialect interface for generate DROP statement
func (pb PostgresDialect) Drop(t *xql.Table, force bool) (stm string, args []interface{}, err error) {
    if nil == t {
        err = errors.New("Table can not be nil !")
        return
    }
    statements := []string{}
    indexes := []*xql.Index{}
    schema := ""
    forced := ""
    if force {
        forced = " IF EXISTS "
    }
    for _, col := range t.GetColumns() {
        indexes = append(indexes, col.Indexes...)
    }
    indexes = append(indexes, t.GetIndexes()...)
    for _, idx := range indexes {
        statements = append(statements, fmt.Sprintf("DROP INDEX %s%s%s;", forced, schema, idx.Name))
    }
    statements = append(statements, fmt.Sprintf("DROP TABLE %s%s;", forced, t.TableName()))
    stm = strings.Join(statements, "\n")
    return
}

// Create
// Implement the IDialect interface for creating table.
func (pb PostgresDialect) Create(t *xql.Table, options ...interface{}) (s string, args []interface{}, err error) {
    var createSQL string
    var table_name string = t.TableName()
    createSQL = "CREATE TABLE IF NOT EXISTS " + table_name + " ( "
    var indexes []*xql.Index
    var cols []string
    for _, c := range t.GetColumns() {
        col_str := fmt.Sprintf(`"%s" %s`, c.FieldName, c.TypeDefine)
        if c.Default != nil {
            col_str = fmt.Sprintf(`%s DEFAULT %s`, col_str, c.Default)
        }
        if len(c.Constraints) > 0 {
            col_str = fmt.Sprintf(`%s %s`, col_str, makeInlineConstraint(c.Constraints...))
        }
        cols = append(cols, col_str)
        indexes = append(indexes, c.Indexes...)
    }
    cols = append(cols, makeConstraints(t, 0, t.GetConstraints()...)...)
    createSQL = createSQL + strings.Join(cols, ", ") + " );"
    indexes = append(indexes, t.GetIndexes()...)
    indexes_strings := makeIndexes(t, 0, indexes...)
    s = strings.Join(append([]string{createSQL}, indexes_strings...), "\n")
    return s, args, err
}

// Select
// Implement the IDialect interface for select values.
func (pb PostgresDialect) Select(t *xql.Table, cols []xql.QueryColumn, filters []xql.QueryFilter, orders []xql.QueryOrder, offset int64, limit int64) (s string, args []interface{}, err error) {
    var colnames []string
    for _, x := range cols {
        colnames = append(colnames, x.String())
    }
    s = fmt.Sprintf("SELECT %s FROM ", strings.Join(colnames, ","))
    s += t.TableName()
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
            } else {
                s = fmt.Sprintf(`%s %s $%d %s %s`, s, cause, n, f.Operator, f.Field)
            }

            args = append(args, f.Value)
        } else {
            n += 1
            if f.Function != "" {
                s = fmt.Sprintf(`%s %s %s %s %s($%d)`, s, cause, f.Field, f.Operator, f.Function, n)
            } else {
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

func isEmptyValue(v reflect.Value) bool {
    switch v.Kind() {
    case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
        return v.Len() == 0
    case reflect.Bool:
        return !v.Bool()
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return v.Int() == 0
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        return v.Uint() == 0
    case reflect.Float32, reflect.Float64:
        return v.Float() == 0
    case reflect.Interface, reflect.Ptr:
        return v.IsNil()
    }
    return false
}

// Insert
// Implement the IDialect interface to generate insert statement
func (pb PostgresDialect) Insert(t *xql.Table, obj interface{}, col ...string) (s string, args []interface{}, err error) {
    s = "INSERT INTO "
    s += t.TableName()
    var cols []string
    var vals []string
    r := reflect.ValueOf(obj)
    if len(col) < 1 {
        for _, x := range t.GetColumns() {
            col = append(col, x.FieldName)
        }
    }
    var i int
    for _, n := range col {
        column, ok := t.GetColumn(n)
        if ! ok {
            continue
        }
        fv := reflect.Indirect(r).FieldByName(column.ElemName)
        fv.Kind()
        //if fv.Interface() == reflect.Zero(fv.Type()).Interface() {
        if ! fv.IsValid() || isEmptyValue(fv) {
        //if ( fv.Kind() == reflect.Ptr && fv.IsNil() ) || reflect.Zero(fv.Type()).Interface() == fv.Interface() {
            if column.Default == nil {
                continue
            }
            args = append(args, column.Default)
        }else{
            args = append(args, fv.Interface())
        }
        i += 1
        cols = append(cols, fmt.Sprintf(`"%s"`, column.FieldName))
        vals = append(vals, fmt.Sprintf("$%d", i))

    }
    s = fmt.Sprintf("%s (%s) VALUES(%s)", s, strings.Join(cols, ","), strings.Join(vals, ","))
    return
}

// Update
// Implement the IDialect interface to generate UPDATE statement
func (pb PostgresDialect) Update(t *xql.Table, filters []xql.QueryFilter, cols ...xql.UpdateColumn) (s string, args []interface{}, err error) {
    s = "UPDATE "
    s += t.TableName()
    if len(cols) < 1 {
        panic("Empty Update MappedColumns!!!")
    }
    var n int
    for i, uc := range cols {
        n += 1
        if i == 0 {
            //s = s + " SET "
            s = fmt.Sprintf(`%s SET "%s"=$%d`, s, uc.Field, n)
        } else {
            s = fmt.Sprintf(`%s, "%s"=$%d`, s, uc.Field, n)
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

// Delete
// Implement the IDialect interface to generate DELETE statement
func (pb PostgresDialect) Delete(t *xql.Table, filters []xql.QueryFilter) (s string, args []interface{}, err error) {
    s = "DELETE FROM "
    s += t.TableName()
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

func CreateSchema(db *sql.DB, schema string) error {
    s := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s",schema)
    //fmt.Println(">>>", s)
    if _, e := db.Exec(s); nil != e {
        return e
    }
    return nil
}

func Initialize_HSTORE(db *sql.DB, schema... string) error {
    s := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS hstore SCHEMA %s",schema[0])
    //fmt.Println(">>>", s)
    if _, e := db.Exec(s); nil != e {
        return e
    }
    return nil
}

func Initialize_UUID(db *sql.DB, schema... string) error {
    s := fmt.Sprintf("CREATE EXTENSION  IF NOT EXISTS \"uuid-ossp\" SCHEMA %s",schema[0])
    //fmt.Println(">>>", s)
    if _, e := db.Exec(s); nil != e {
        return e
    }
    return nil
}

// Register the dialect.
func init() {
    xql.RegisterDialect("postgres", &PostgresDialect{})
}
