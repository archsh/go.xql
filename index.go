package xql

import "fmt"

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


func makeIndexes(t uint8, name string, fields ...interface{}) []*Index {
    var indexes []*Index
    if t >= INDEX_INVALID || t <= INDEX_NONE {
        return indexes
    }
    for i, f := range fields {
        if nil == f {
            continue
        }
        idx := &Index{Type:t, Name:fmt.Sprintf("%s_%d", name, i+1)}
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