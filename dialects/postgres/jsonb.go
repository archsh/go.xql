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
	return JsonbScan(j, value)
}

func (j JSON) Value() (driver.Value, error) {
	return JsonbValue(j)
}

type JSONB JSON

func (j JSONB) Declare(props xql.PropertySet) string {
	return "JSONB"
}

func (j *JSONB) Scan(value interface{}) error {
	return JsonbScan(j, value)
}

func (j JSONB) Value() (driver.Value, error) {
	return JsonbValue(j)
}

func JsonbScan(dest interface{}, src interface{}) error {
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

func JsonbValue(obj interface{}) (driver.Value, error) {
	return json.Marshal(obj)
}
