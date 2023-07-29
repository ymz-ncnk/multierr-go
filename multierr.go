package multierr

import (
	"reflect"
	"sort"
	"strings"
)

// Sep separates one error message from another.
const Sep = "; "

// New creates a new multiError.
func New(es []error) error {
	if len(es) == 0 {
		return nil
	}
	return &multiError{es}
}

// multiError designed to combine several errors into one.
type multiError struct {
	es []error
}

// Get returns the i-th error.
func (e *multiError) Get(i int) error {
	return e.es[i]
}

// Len returns the number of errors.
func (e *multiError) Len() int {
	return len(e.es)
}

// Unwrap makes a new slice from the existing errors.
func (e *multiError) Unwrap() []error {
	es := make([]error, len(e.es))
	copy(es, e.es)
	return es
}

func (e *multiError) Error() string {
	n := len(Sep) * (len(e.es) - 1)
	for i := 0; i < len(e.es); i++ {
		n += len(e.es[i].Error())
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(e.es[0].Error())
	for _, err := range e.es[1:] {
		b.WriteString(Sep)
		b.WriteString(err.Error())
	}
	return b.String()
}

// Similar checks if two multiErrors are similar. Two multiErrors are
// considered similar if they contain the same errors, even in a different
// order.
func (e *multiError) Similar(ae error) bool {
	if ae == nil {
		return false
	}
	mae, ok := ae.(*multiError)
	if !ok {
		return false
	}

	if e.Len() != mae.Len() {
		return false
	}
	var (
		es1 = sortErrors(e.es)
		es2 = sortErrors(mae.es)
	)
	for i := 0; i < len(es1); i++ {
		if !similarErrors(es1[i], es2[i]) {
			return false
		}
	}
	return true
}

func sortErrors(es []error) []error {
	a := make([]error, len(es))
	copy(a, es)
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].Error() < a[j].Error()
	})
	return a
}

func similarErrors(e1, e2 error) bool {
	if reflect.TypeOf(e1) != reflect.TypeOf(e2) {
		return false
	}
	return e1.Error() == e2.Error()
}
