package dbconn

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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

