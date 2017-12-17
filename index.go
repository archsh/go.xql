package xql

type Index struct {
    Fields []*Column
    Type   string
}


func makeIndexes(fields ...interface{}) []*Index {
    var indexes []*Index
    for _, f := range fields {
        idx := &Index{}
        if fc, ok := f.(*Column); ok {
            idx.Fields = []*Column{fc}

        }else if fcs, ok := f.([]*Column); ok {
            idx.Fields = fcs
        }else if fs, ok := f.(string); ok {

        }else{
            continue
        }
        indexes = append(indexes, idx)
    }
    for _,field := range fields {
        if nil == field {
            continue
        }
        if field.Indexed {
            idx := &Index{Field:field}
            indexes = append(indexes, idx)
        }else{
            continue
        }
    }
    return indexes
}