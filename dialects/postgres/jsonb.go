package postgres

import (
    "errors"
    "reflect"
    "database/sql/driver"
    "encoding/json"
)

func JSONB_Scan(dest interface{}, src interface{}) error {
    if nil == src {
        dest = nil
        return nil
    }
    source, ok := src.([]byte)
    if !ok {
        return errors.New("Type assertion .([]byte) failed.")
    }
    entityType := reflect.TypeOf(dest)
    obj := reflect.New(entityType)
    err := json.Unmarshal(source, obj.Elem().Addr().Interface())
    if err != nil {
        return err
    }
    //reflect.Indirect(dest).Set(obj)
    dest = obj.Elem().Addr().Interface()
    return nil
}

func JSONB_Value(obj interface{}) (driver.Value, error) {
    return json.Marshal(obj)
}
