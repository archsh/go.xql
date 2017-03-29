package postgres
import (
    "github.com/archsh/go.xql"
    "fmt"
    "reflect"
    "strings"
    "errors"
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

func (pb PostgresDialect) Create(t *xql.Table, options ...interface{}) (s string, args []interface{}, err error) {
    var createSQL string
    if t.Schema != "" {
        createSQL = "CREATE TABLE IF NOT EXISTS " + t.Schema+"."+t.TableName + " ( "
    }else{
        createSQL = "CREATE TABLE IF NOT EXISTS " + t.TableName + " ( "
    }
    var indexes []string
    var cols []string
    for _, c := range t.ListedColumns {
        if c.Type == "" {
            return "", args, errors.New("Unknow Column Type!!!")
        }
        col_str := fmt.Sprintf(`"%s" %s`, c.FieldName, c.Type)
        if ! c.Nullable {
            col_str = fmt.Sprintf(`%s NOT NULL`, col_str)
        }
        cols = append(cols, col_str)
        if c.Indexed {
            idxs := fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_%s ON %s USING btree (%s);",
                t.TableName, c.FieldName, t.TableName, c.FieldName)
            indexes = append(indexes, idxs)
        }
    }
    if len(t.PrimaryKey) > 0 {
        cols = append(cols, fmt.Sprintf("CONSTRAINT %s_pkey PRIMARY KEY (%s)", t.TableName, strings.Join(t.PrimaryKey, ",")))
    }
    createSQL =  createSQL + strings.Join(cols, ",") + " );"
    s = strings.Join(append([]string{createSQL}, indexes...), "\n")
    return s, args, err
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
            v, ok := t.MappedColumns[n]
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
        for k, v := range t.MappedColumns {
            //fmt.Println("POSTGRES Insert>2>>>",k,v.Auto,v.PrimaryKey,v)
            if v.Auto {
                continue
            }
            i += 1
            cols = append(cols, fmt.Sprintf(`"%s"`,k))
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
        panic("Empty Update MappedColumns!!!")
    }
    var n int
    for i, uc := range cols {
        n += 1
        if i == 0 {
            //s = s + " SET "
            s = fmt.Sprintf(`%s SET "%s"=$%d`,s, uc.Field, n)
        }else{
            s = fmt.Sprintf(`%s, "%s"=$%d`,s, uc.Field, n)
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