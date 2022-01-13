package errors

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {
	t.Run("Check that error returned.", func(t *testing.T) {
		if err := InternalServerError.New("test error"); err == nil {
			t.Fatalf("The expected error message does not come back.")
		}
	})
}

func Test_Wrap(t *testing.T) {
	t.Run("Check that error returned.", func(t *testing.T) {
		e := fmt.Errorf("yyy")
		if err := InternalServerError.Wrap("test error", e); err == nil {
			t.Fatalf("The expected error message does not come back.")
		}
	})
}
