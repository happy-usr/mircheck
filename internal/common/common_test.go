package common

import (
	"errors"
	sqlite3 "github.com/mattn/go-sqlite3"
	"testing"
)

func TestOpen(t *testing.T) {
	tests := []struct {
		db          DBSelection
		testName    string
		wantErr     bool
		expectedErr error
		isDBPtrNull bool
	}{
		{TEST, "Test database", false, nil, false},
		{DEFAULT, "Default database", false, nil, false},
		{208, "Invalid database", true, ErrInvalidDBSelection, true},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
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
				t.Fatalf("%s: pointer 'db' has an "+
					"unexpected value\n",
					test.testName)
			}
			if !test.isDBPtrNull {
				err := db.Ping()
				if err != nil {
					t.Fatalf("%s: failed to "+
						"establish a "+
						" database "+
						"connection: %s\n",
						test.testName,
						err.Error())
				}
			}
		})
	}
}

type TestIsErrSqlConstraintStruct struct {
	err            error
	sqlConstraint  sqlite3.ErrNoExtended
	testName       string
	expectedOutput bool
}

type tmpStruct struct {
}

func (ts tmpStruct) Error() string {
	return ""
}

func makeAndAddErrs(tests *[]TestIsErrSqlConstraintStruct) {
	sqlite3Errs := []sqlite3.Error{
		{Code: sqlite3.ErrAuth, /*anything other than ErrConstraint*/
			ExtendedCode: sqlite3.ErrConstraintRowID},
		{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintCheck},
		{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintNotNull},
		{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintPrimaryKey},
		{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintUnique},
		{Code: sqlite3.ErrConstraint, ExtendedCode: sqlite3.ErrConstraintRowID},
	}
	var tmpTest TestIsErrSqlConstraintStruct
	tmpTest.err = tmpStruct{} /*anything other than expected error*/
	tmpTest.sqlConstraint = sqlite3Errs[0].ExtendedCode
	tmpTest.testName = "unexpected error type"
	tmpTest.expectedOutput = false
	*tests = append(*tests, tmpTest)

	tmpTest.err = sqlite3Errs[0]
	tmpTest.sqlConstraint = sqlite3Errs[0].ExtendedCode
	tmpTest.testName = "unexpected error code"
	tmpTest.expectedOutput = false
	*tests = append(*tests, tmpTest)

	tmpTest.err = sqlite3Errs[1]
	tmpTest.sqlConstraint = sqlite3Errs[1].ExtendedCode
	tmpTest.testName = "check constraint"
	tmpTest.expectedOutput = true
	*tests = append(*tests, tmpTest)

	tmpTest.err = sqlite3Errs[2]
	tmpTest.sqlConstraint = sqlite3Errs[2].ExtendedCode
	tmpTest.testName = "NotNull constraint"
	tmpTest.expectedOutput = true
	*tests = append(*tests, tmpTest)

	tmpTest.err = sqlite3Errs[3]
	tmpTest.sqlConstraint = sqlite3Errs[4].ExtendedCode
	tmpTest.testName = "PK constraint"
	tmpTest.expectedOutput = false
	*tests = append(*tests, tmpTest)

	tmpTest.err = sqlite3Errs[4]
	tmpTest.sqlConstraint = sqlite3Errs[5].ExtendedCode
	tmpTest.testName = "unique constraint"
	tmpTest.expectedOutput = false
	*tests = append(*tests, tmpTest)
}

func TestIsErrSqliteConstraint(t *testing.T) {
	tests := []TestIsErrSqlConstraintStruct{
		{nil, 0, "nil error 1", true},
		{nil, sqlite3.ErrIoErrRead /*anything other than ErrConstraints*/, "nil error 2", false},
	}
	makeAndAddErrs(&tests)
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			t.Parallel()
			got := IsErrSqliteConstraint(test.err,
				test.sqlConstraint)
			if test.expectedOutput != got {
				t.Fatalf("%s: expected %t, got %t\n",
					test.testName, test.expectedOutput,
					got)
			}
		})
	}
}
