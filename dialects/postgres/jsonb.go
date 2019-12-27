package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/archsh/go.xql"
)

type JSON struct {
}

func (j JSON) Declare(props xql.PropertySet) string {
	return "json"
}

func (j *JSON) Scan(value interface{}) error {
	return JSONB_Scan(j, value)
}

func (j JSON) Value() (driver.Value, error) {
	return JSONB_Value(j)
}

type JSONB JSON

func (j JSONB) Declare(props xql.PropertySet) string {
	return "JSONB"
}

func (j *JSONB) Scan(value interface{}) error {
	return JSONB_Scan(j, value)
}

func (j JSONB) Value() (driver.Value, error) {
	return JSONB_Value(j)
}

func JSONB_Scan(dest interface{}, src interface{}) error {
	if nil == src {
		dest = nil
		return nil
	}
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}
	//entityType := reflect.TypeOf(dest)
	//obj := reflect.New(entityType)
	err := json.Unmarshal(source, dest)
	if err != nil {
		return err
	}
	//reflect.Indirect(dest).Set(obj)
	//dest = obj.Elem().Addr().Interface()
	return nil
}

func JSONB_Value(obj interface{}) (driver.Value, error) {
	return json.Marshal(obj)
}
