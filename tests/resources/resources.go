package resources

import (
	"testing"
)

type UnitTestSuite[T any, U any] []UnitTestCase[T, U]

type UnitTestCase[T any, U any] struct {
	It     string
	Before func(t *testing.T, in *T)
	In     T
	Want   U
}
