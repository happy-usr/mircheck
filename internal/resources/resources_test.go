package resources

import (
	"database/sql"
	"fmt"
	"github.com/happy-usr/mircheck/internal/common"
	sqlite3 "github.com/mattn/go-sqlite3"
	"testing"
)

type TestResource struct {
	resource    Resource
	testName    string
	wantErr     bool
	expectedErr sqlite3.ErrNoExtended
}

var gResource Resource
var gDB *sql.DB

func init() {
	gDB, _ = common.Open(common.TEST)
	gResource = Resource{gDB, "GLOBAL_T_1", "GLOBAL_R_2"}
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
		{gResource, "simple test", false, 0},
		{gResource, "PK duplication", true, sqlite3.ErrConstraintPrimaryKey},
		{Resource{gDB, "ADD_T_1", ""}, "null PK", true, sqlite3.ErrConstraintCheck},
		{Resource{gDB, "", "ADD_R_1"}, "null 'type' column", true, sqlite3.ErrConstraintCheck},
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
	tests := []TestResource{
		{gResource, "simple test", false, 0},
	}

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

func TestGetResources(t *testing.T) {
	testName := "simple test"
	resourcesPtr, err := GetResources(gDB)
	if resourcesPtr == nil {
		t.Fatalf("%s: pointer to resources is nil\n", testName)
	}
	if err != nil {
		t.Fatalf("%s: unexpected error: %s\n", testName, err.Error())
	}
}

type TestUpdate struct {
	resource    Resource
	updateVal   string
	testName    string
	wantErr     bool
	expectedErr sqlite3.ErrNoExtended
}

func TestUpdateResource(t *testing.T) {
	defer gDB.Close()

	resource1 := Resource{gDB, "UPDATE_T_1", "UPDATE_R_1"}
	resource2 := Resource{gDB, "UPDATE_T_2", "UPDATE_R_2"}
	resource1UpdateVal := "UPDATE_CHANGED_R_1"
	updatedResource1 := Resource{gDB, "UPDATE_T_1", resource1UpdateVal}
	AddResource(resource1)
	AddResource(resource2)

	tests := []TestUpdate{
		{resource1, resource1UpdateVal, "simple test", false, 0},
		{updatedResource1, resource2.Resource, "PK test", true,
			sqlite3.ErrConstraintPrimaryKey},
		{updatedResource1, "", "null PK test", true,
			sqlite3.ErrConstraintCheck},
	}
	for _, test := range tests {
		err := UpdateResource(test.resource, test.updateVal)
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
