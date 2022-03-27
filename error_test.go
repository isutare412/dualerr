package dualerr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/isutare412/dualerr"
)

func TestWrapping(t *testing.T) {
	var (
		systemErrMsg  = "base system error"
		simpleErrMsg  = "client not found"
		wrappedErrMsg = fmt.Sprintf("dualerr_test.TestWrapping: %s", systemErrMsg)
	)

	base := fmt.Errorf(systemErrMsg)
	wrap := error(dualerr.New(base, simpleErrMsg))

	if !errors.Is(wrap, dualerr.Err) {
		t.Fatalf("error is not DualErr")
	}

	if wrap.Error() != wrappedErrMsg {
		t.Fatalf("error message is not equal: got[%s] expected[%s]",
			wrap.Error(), wrappedErrMsg)
	}

	var restored dualerr.Error
	if !errors.As(wrap, &restored) {
		t.Fatalf("cannot cast error to DualError")
	}

	if restored.Error() != wrappedErrMsg {
		t.Fatalf("error message is not equal: got[%s] expected[%s]",
			restored.Error(), wrappedErrMsg)
	}
	if restored.Simple().Error() != simpleErrMsg {
		t.Fatalf("error simple message is not equal: got[%s] expected[%s]",
			restored.Simple(), simpleErrMsg)
	}
}
