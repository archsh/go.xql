package xql

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
    Type uint8
    Fields []*Column
}


func makeIndexes(t uint8, fields ...interface{}) []*Index {
    var indexes []*Index
    if t >= INDEX_INVALID || t <= INDEX_NONE {
        return indexes
    }
    for _, f := range fields {
        if nil == f {
            continue
        }
        idx := &Index{Type:t}
        if fc, ok := f.(*Column); ok {
            idx.Fields = []*Column{fc}
        }else if fcs, ok := f.([]*Column); ok {
            idx.Fields = fcs
        }/*else if fs, ok := f.(string); ok {

        }*/else{
            continue
        }
        indexes = append(indexes, idx)
    }

    return indexes
}