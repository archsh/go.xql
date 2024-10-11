package postgres

import (
	"database/sql"
	"database/sql/driver"
	xql "github.com/archsh/go.xql"
)

type UUID string

func (j UUID) Declare(props xql.PropertySet) string {
	return "uuid"
}

func (j *UUID) Scan(value interface{}) error {
	var v sql.NullString
	if err := v.Scan(value); err != nil {
		return err
	} else if v.Valid {
		*j = UUID(v.String)
	} else {
		*j = UUID("")
	}
	return nil
}

func (j UUID) Value() (driver.Value, error) {
	if j == "" {
		return nil, nil
	} else {
		return j, nil
	}
}

func (j UUID) String() string {
	return string(j)
}
