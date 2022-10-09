package resources

import "testing"

type UnitTestSuite[T any, U any] []UnitTestCase[T, U]

type UnitTestCase[T any, U any] struct {
	It     string
	Before func(t *testing.T, in *T)
	In     T
	Want   U
}

func (uts UnitTestSuite[T, U]) Run(t *testing.T, fn func(t *testing.T, testCase UnitTestCase[T, U])) {
	for i := range uts {
		i := i
		tc := uts[i]
		t.Run(tc.It, func(t *testing.T) {
			t.Parallel()
			fn(t, tc)
		})
		uts[i] = tc
	}
}
