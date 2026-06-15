package common

import (
	"database/sql"
	sqlite3 "github.com/mattn/go-sqlite3"
	"errors"
)

type DBSelection int
const (
	TEST DBSelection = iota
	DEFAULT
	NONE
)
var ErrInvalidDBSelection = errors.New("invalid database")

func Open(db DBSelection) (*sql.DB, error) {
	var dbAddr string
	switch db {
		case TEST:
			dbAddr = "../../test.db"
		case DEFAULT:
			dbAddr = "../../mirrors.db"
		default:
			return nil, ErrInvalidDBSelection
	}
	return sql.Open("sqlite3", dbAddr)
}

func CmpDriverWithErrSqliteConstraint(err error, 
	sqliteConstraint sqlite3.ErrNoExtended) bool {
	if err == nil && sqliteConstraint == 0 {
		return true
	}
	var tmpSqliteErr sqlite3.Error
	cmp := errors.As(err, &tmpSqliteErr)
	if !cmp {
		return false
	}
	if tmpSqliteErr.Code != sqlite3.ErrConstraint ||
		tmpSqliteErr.ExtendedCode != sqliteConstraint {
		return false
	}
	return true
}

