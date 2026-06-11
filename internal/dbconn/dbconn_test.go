package dbconn

import (
	"testing"
	"errors"
)

type dbConnTest struct {
	db DBSelection
	testName string
	wantErr bool
	expectedErr error
	isDBPtrNull bool
}

func TestOpen(t *testing.T) {
	tests := []dbConnTest {
		{TEST, "Test database", false, nil, false},
		{DEFAULT, "Default database", false, nil, false},
		{208, "Invalid database", true, ErrInvalidDBSelection, true},
	}
	for _, test := range tests {
		t.Run(test.testName, func (t *testing.T) {
			t.Parallel()
			db, err := Open(test.db)
			if test.wantErr != (err != nil) {
				t.Fatalf(`%s: unexpected error value; 
					%s\n`, test.testName,
					err.Error())
			}
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("%s: unexpected error type\n",
					test.testName)
			}
			if test.isDBPtrNull != (db == nil) {
				t.Fatalf(`%s: pointer 'db' has an 
					unexpected value\n`, test.testName)
			}
			if !test.isDBPtrNull {
				err := db.Ping()
				if err != nil {
					t.Fatalf(`%s: failed to establish a 
						database connection; %s\n`,
						test.testName, err.Error())
				}
			}
		})
	}
}
