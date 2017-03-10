package xql

import "database/sql"

func MakeSession(db *sql.DB, driverName string, verbose ...bool) *Session {
    sess := &Session{db:db, driverName:driverName}
    if len(verbose) > 0 {
        sess.verbose = verbose[0]
    }
    return sess
}



