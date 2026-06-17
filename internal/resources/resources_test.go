package resources

import (
	"database/sql"
	"fmt"
	"github.com/happy-usr/mircheck/internal/common"
	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/xyproto/randomstring"
	"testing"
)

type TestResource struct {
	resource    Resource
	testName    string
	wantErr     bool
	expectedErr sqlite3.ErrNoExtended
}

var gType, gResource string
var gDB *sql.DB

func init() {
	gType, gResource = genRandomString(), genRandomString()
	gDB, _ = common.Open(common.TEST)
}

func genRandomString() string {
	return randomstring.HumanFriendlyEnglishString(10)
}

func diagErr(err error) string {
	var got_val string
	if err == nil {
		got_val = "(error is nil)"
	} else {
		got_val = fmt.Sprintf("error: %q", err.Error())
	}
	return got_val
}

func TestAddResource(t *testing.T) {
	tests := []TestResource{
		{Resource{gDB, gType, gResource}, "simple test", false, 0},
		{Resource{gDB, gType, gResource}, "PK duplication", true, sqlite3.ErrConstraintPrimaryKey},
		{Resource{gDB, gType, ""}, "null PK", true, sqlite3.ErrConstraintCheck},
		{Resource{gDB, "", genRandomString()}, "null 'type' column", true, sqlite3.ErrConstraintCheck},
	}

	for _, test := range tests {
		err := AddResource(test.resource)
		if test.wantErr != (err != nil) {
			got_val := diagErr(err)
			t.Fatalf("%s: unexpected error value; %s\n",
				test.testName, got_val)
		}
		if test.wantErr {
			if !common.IsErrSqliteConstraint(err, test.expectedErr) {
				t.Fatalf("%s: unexpected error type; expected "+
					"%q, got %q\n", test.testName,
					test.expectedErr.Error(), err.Error())
			}
		}
	}
}

func TestRemoveResource(t *testing.T) {
	tests := []TestResource {
		{Resource{gDB, gType, gResource}, "simple test", false, 0},
	}
	defer gDB.Close()

	for _, test := range tests {
		err := RemoveResource(test.resource)
        if test.wantErr != (err != nil) {
            got_val := diagErr(err)
            t.Fatalf("%s: unexpected error value; %s\n",
                test.testName, got_val)
        }
        if test.wantErr {
            if !common.IsErrSqliteConstraint(err, test.expectedErr) {
                t.Fatalf("%s: unexpected error type; expected "+
                    "%q, got %q\n", test.testName,
                    test.expectedErr.Error(), err.Error())
            }
        }
	}
}

