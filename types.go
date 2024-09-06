package xql

import (
	"database/sql"
	"encoding/json"
)

type Field[T any] sql.Null[T]

func (f Field[T]) UnmarshalJSON(bytes []byte) error {
	//TODO implement me
	//panic("implement me")
	var t = new(T)
	if e := json.Unmarshal(bytes, t); nil != e {
		return e
	} else {
		f.V, f.Valid = *t, true
		return nil
	}
}

func (f Field[T]) MarshalJSON() ([]byte, error) {
	//TODO implement me
	//panic("implement me")
	if f.Valid {
		return json.Marshal(f.V)
	} else {
		return json.Marshal(nil)
	}
}
