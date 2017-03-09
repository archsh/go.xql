package xql

import "database/sql"

func MakeSession(db *sql.DB, driverName string) *Session {
    return &Session{db:db, driverName:driverName}
}



