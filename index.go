package xql

import (
    "fmt"
    "strings"
)

// PostgreSQL provides several index types: B-tree, Hash, GiST, SP-GiST, GIN and BRIN
const (
    INDEX_NONE uint8 = iota
    INDEX_B_TREE
    INDEX_HASH
    INDEX_GIST
    INDEX_SP_GIST
    INDEX_GIN
    INDEX_BRIN
    INDEX_INVALID
)

type Index struct {
    Type    uint8
    Name    string
    Columns []*Column
}

func buildIndexes(t *Table, ss... [2]string) []*Index {
    var indexes []*Index
    for _,xs := range ss {
        idx := &Index{}
        switch strings.ToLower(xs[0]) {
        case "hash":
            idx.Type = INDEX_HASH
        case "gist":
            idx.Type = INDEX_GIST
        case "sp_gist","sp-gist":
            idx.Type = INDEX_SP_GIST
        case "brin":
            idx.Type = INDEX_BRIN
        case "gin":
            idx.Type = INDEX_GIN
        default:
            idx.Type = INDEX_B_TREE
        }
        for _,f := range strings.Split(xs[1],",") {
            if field, ok := t.GetColumn(f); ok {
                idx.Columns = append(idx.Columns, field)
            }else{
                panic(fmt.Sprintf("Can not get column '%s' from '%s'!", f, t.TableName()))
            }
        }
        idx.Name = fmt.Sprintf("%s_%s_idx", t.BaseTableName(), strings.Join(strings.Split(xs[1],","),"_"))
        indexes = append(indexes, idx)
    }
    return indexes
}


func makeIndexes(t uint8, name string, fields ...interface{}) []*Index {
    var indexes []*Index
    if t >= INDEX_INVALID || t <= INDEX_NONE {
        return indexes
    }
    for i, f := range fields {
        if nil == f {
            continue
        }
        idx := &Index{Type:t, Name:fmt.Sprintf("%s_%d_idx", name, i+1)}
        if fc, ok := f.(*Column); ok {
            idx.Columns = []*Column{fc}
        }else if fcs, ok := f.([]*Column); ok {
            idx.Columns = fcs
        }else{
            continue
        }
        indexes = append(indexes, idx)
    }
    return indexes
}