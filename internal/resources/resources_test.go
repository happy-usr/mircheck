package resources

import (
	"database/sql"
	"github.com/happy-usr/mircheck/internal/common"
	"github.com/happy-usr/mircheck/internal/commontest"
	sqlite3 "github.com/mattn/go-sqlite3"
	"testing"
)

type TestResource struct {
	resource Resource
	err      commontest.TestSqliteConstraint
}

var gResource Resource
var gDB *sql.DB

func init() {
	gDB, _ = common.Open(common.TEST)
	gResource = Resource{gDB, "GLOBAL_T_1", "GLOBAL_R_2"}
}

func TestAddResource(t *testing.T) {
	tests := []TestResource{
		{gResource, commontest.TestSqliteConstraint{
			WantErr:        false,
			ExpectedSqlErr: 0,
			TestName:       "simple test",
		},
		},
		{gResource, commontest.TestSqliteConstraint{
			WantErr:        true,
			ExpectedSqlErr: sqlite3.ErrConstraintPrimaryKey,
			TestName:       "PK duplication",
		},
		},
		{Resource{gDB, "ADD_T_1", ""}, commontest.TestSqliteConstraint{
			WantErr:        true,
			ExpectedSqlErr: sqlite3.ErrConstraintCheck,
			TestName:       "null PK",
		},
		},
		{Resource{gDB, "", "ADD_R_1"}, commontest.TestSqliteConstraint{
			WantErr:        true,
			ExpectedSqlErr: sqlite3.ErrConstraintCheck,
			TestName:       "null 'type' column",
		},
		},
	}

	for _, test := range tests {
		test.err.Err = AddResource(test.resource)
		errString := test.err.CheckErrSqlConstraint()
		if errString != "" {
			t.Fatal(errString)
		}
	}
}

func TestRemoveResource(t *testing.T) {
	tests := []TestResource{
		{gResource, commontest.TestSqliteConstraint{
			WantErr:        false,
			ExpectedSqlErr: 0,
			TestName:       "simple test",
		},
		},
	}

	for _, test := range tests {
		test.err.Err = RemoveResource(test.resource)
		errString := test.err.CheckErrSqlConstraint()
		if errString != "" {
			t.Fatal(errString)
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
	resource  Resource
	updateVal string
	err       commontest.TestSqliteConstraint
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
		{resource1, resource1UpdateVal, commontest.TestSqliteConstraint{
			WantErr:        false,
			ExpectedSqlErr: 0,
			TestName:       "simple test",
		},
		},
		{updatedResource1, resource2.Resource, commontest.TestSqliteConstraint{
			WantErr:        true,
			ExpectedSqlErr: sqlite3.ErrConstraintPrimaryKey,
			TestName:       "PK test",
		},
		},
		{updatedResource1, "", commontest.TestSqliteConstraint{
			WantErr:        true,
			ExpectedSqlErr: sqlite3.ErrConstraintCheck,
			TestName:       "null PK test",
		},
		},
	}
	for _, test := range tests {
		test.err.Err = UpdateResource(test.resource, test.updateVal)
		errString := test.err.CheckErrSqlConstraint()
		if errString != "" {
			t.Fatal(errString)
		}
	}
}
