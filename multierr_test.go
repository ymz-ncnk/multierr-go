package multierr

import (
	"errors"
	"testing"
)

type anotherError struct{}

func (e anotherError) Error() string {
	return ""
}

func TestMultiError(t *testing.T) {

	t.Run("New should create a correct error", func(t *testing.T) {
		var (
			wantErr1 = errors.New("error1")
			wantErr2 = errors.New("error2")
			wantStr  = wantErr1.Error() + Sep + wantErr2.Error()
			wantLen  = 2
			es       = []error{wantErr1, wantErr2}
			err      = New(es)
			merr     = err.(*multiError)
		)
		if err.Error() != wantStr {
			t.Errorf("unexpected str, want '%v' actual '%v'", wantStr, err.Error())
		}
		if merr.Len() != 2 {
			t.Errorf("unexpected len, want '%v' actual '%v'", wantLen, merr.Len())
		}
		if merr.Get(0) != wantErr1 {
			t.Errorf("unexpected first error, want '%v' actual '%v'", wantErr1,
				merr.Get(0))
		}
		if merr.Get(1) != wantErr2 {
			t.Errorf("unexpected second error, want '%v' actual '%v'", wantErr2,
				merr.Get(1))
		}
	})

	t.Run("If New receives an empty slice it should return nil",
		func(t *testing.T) {
			err := New([]error{})
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
		})

	t.Run("Two multi errors with different error orders must be similar",
		func(t *testing.T) {
			var (
				merr1 = New([]error{
					errors.New("error1"), errors.New("error2"),
				}).(*multiError)
				merr2 = New([]error{
					errors.New("error2"), errors.New("error1"),
				}).(*multiError)
			)
			if !merr1.Similar(merr2) {
				t.Error("errors are not similar")
			}
		})

	t.Run("Different multierrors should be unsimilar", func(t *testing.T) {
		var (
			merr1 = New([]error{
				errors.New("error1"), errors.New("error2"),
			}).(*multiError)
			merr2 = New([]error{
				errors.New("error3"), errors.New("error2"),
			}).(*multiError)
		)
		if merr1.Similar(merr2) {
			t.Error("errors should be unsimilar")
		}
	})

	t.Run("Unwrap should return know slice", func(t *testing.T) {
		var (
			wantES = []error{errors.New("error1"), errors.New("error2")}
			merr   = New(wantES)
			es     = merr.(*multiError).Unwrap()
		)
		if len(es) != len(wantES) {
			t.Error("Unwarp returns wron slice")
		}
		for i := 0; i < len(es); i++ {
			if !similarErrors(es[i], wantES[i]) {
				t.Error("Unwarp returns wron slice")
			}
		}
	})

	t.Run("Similar method should return false, if it received nil",
		func(t *testing.T) {
			merr := New([]error{errors.New("error1")})
			if merr.(*multiError).Similar(nil) {
				t.Error("Similar method does not work correctly")
			}
		})

	t.Run("Similar method should return false, if it did not receive a multiError",
		func(t *testing.T) {
			merr := New([]error{errors.New("error1")})
			if merr.(*multiError).Similar(errors.New("error2")) {
				t.Error("Similar method does not work correctly")
			}
		})

	t.Run("Similar method should return false, if it did not receive a multiError with the same length",
		func(t *testing.T) {
			var (
				merr1 = New([]error{errors.New("error1")})
				merr2 = New([]error{errors.New("error1"), errors.New("error2")})
			)
			if merr1.(*multiError).Similar(merr2) {
				t.Error("Similar method does not work correctly")
			}
		})

	t.Run("Similar method should return false, if one multiError contains an error of different type",
		func(t *testing.T) {
			var (
				merr1 = New([]error{anotherError{}})
				merr2 = New([]error{errors.New("error1")})
			)
			if merr1.(*multiError).Similar(merr2) {
				t.Error("Similar method does not work correctly")
			}
		})

}
