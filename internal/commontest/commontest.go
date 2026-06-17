package commontest

import (
	"fmt"
	"github.com/happy-usr/mircheck/internal/common"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type TestSqliteConstraint struct {
	Err            error
	WantErr        bool
	ExpectedSqlErr sqlite3.ErrNoExtended
	TestName       string
}

func diagErr(err error) string {
	var got string
	if err == nil {
		got = "(error is nil)"
	} else {
		got = fmt.Sprintf("error: %q", err.Error())
	}
	return got
}

func (s *TestSqliteConstraint) CheckErrSqlConstraint() string {
	if s.WantErr != (s.Err != nil) {
		got := diagErr(s.Err)
		return fmt.Sprintf("%s: unexpected error value; %s\n",
			s.TestName, got)
	}
	if s.WantErr {
		if !common.IsErrSqliteConstraint(s.Err, s.ExpectedSqlErr) {
			fmt.Sprintf("%s: unexpected error type; expected "+
				"%q, got %q\n", s.TestName,
				s.ExpectedSqlErr.Error(), s.Err.Error())
		}
	}
	return ""
}
