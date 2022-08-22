package assert

import (
	"fmt"
	"testing"
)

type Data struct {
	Str string
	Num int64
}

func TestEqual(t *testing.T) {
	assertOk(t, "IdenticalStruct", func(t testing.TB) {
		Equal(t, Data{"expected", 1234}, Data{"expected", 1234})
	})
	assertFail(t, "DifferentStruct", func(t testing.TB) {
		Equal(t, Data{"expected\ntext", 1234}, Data{"actual\ntext", 1234})
	})
}

func TestEqualStrings(t *testing.T) {
	assertFail(t, "IdenticalStrings", func(t testing.TB) {
		Equal(t, "hello\nworld", "goodbye\nworld")
	})
}

func TestNotEqual(t *testing.T) {
	assertOk(t, "DifferentFieldValue", func(t testing.TB) {
		NotEqual(t, Data{"expected", 1234}, Data{"expected", 1235})
	})
	assertFail(t, "SameValue", func(t testing.TB) {
		NotEqual(t, Data{"expected", 1234}, Data{"expected", 1234})
	})
}

func TestContains(t *testing.T) {
	assertOk(t, "Found", func(t testing.TB) {
		Contains(t, "a haystack with a needle in it", "needle")
	})
	assertFail(t, "NotFound", func(t testing.TB) {
		Contains(t, "a haystack with a needle in it", "screw")
	})
}

func TestContainsItem(t *testing.T) {
	assertOk(t, "Found", func(t testing.TB) {
		ContainsItem(t, []string{"hello", "world"}, "hello")
	})
	assertFail(t, "NotFound", func(t testing.TB) {
		ContainsItem(t, []string{"hello", "world"}, "goodbye")
	})
}

func TestNotContains(t *testing.T) {
	assertOk(t, "NotFound", func(t testing.TB) {
		NotContains(t, "a haystack with a needle in it", "screw")
	})
	assertFail(t, "Found", func(t testing.TB) {
		NotContains(t, "a haystack with a needle in it", "needle")
	})
}

func TestEqualError(t *testing.T) {
	assertOk(t, "SameMessage", func(t testing.TB) {
		EqualError(t, fmt.Errorf("hello"), "hello")
	})
	assertOk(t, "Nil", func(t testing.TB) {
		EqualError(t, nil, "")
	})
	assertFail(t, "MessageMismatch", func(t testing.TB) {
		EqualError(t, fmt.Errorf("hello"), "goodbye")
	})
}

func TestZero(t *testing.T) {
	assertOk(t, "Struct", func(t testing.TB) {
		Zero(t, Data{})
	})
	assertOk(t, "NilSlice", func(t testing.TB) {
		var slice []int
		Zero(t, slice)
	})
	assertFail(t, "NonEmptyStruct", func(t testing.TB) {
		Zero(t, Data{Str: "str"})
	})
	assertFail(t, "NonEmptySlice", func(t testing.TB) {
		slice := []int{1, 2, 3}
		Zero(t, slice)
	})
	assertOk(t, "ZeroLenSlice", func(t testing.TB) {
		slice := []int{}
		Zero(t, slice)
	})
}

func TestNotZero(t *testing.T) {
	assertOk(t, "PopulatedStruct", func(t testing.TB) {
		notZero := Data{Str: "hello"}
		NotZero(t, notZero)
	})
	assertFail(t, "EmptyStruct", func(t testing.TB) {
		zero := Data{}
		NotZero(t, zero)
	})
	assertFail(t, "NilSlice", func(t testing.TB) {
		var slice []int
		NotZero(t, slice)
	})
	assertFail(t, "ZeroLenSlice", func(t testing.TB) {
		slice := []int{}
		NotZero(t, slice)
	})
	assertOk(t, "Slice", func(t testing.TB) {
		slice := []int{1, 2, 3}
		NotZero(t, slice)
	})
}

type testTester struct {
	*testing.T
	failed string
}

func (t *testTester) Fatalf(message string, args ...interface{}) {
	t.failed = fmt.Sprintf(message, args...)
}

func (t *testTester) Fatal(args ...interface{}) {
	t.failed = fmt.Sprint(args...)
}

func assertFail(t *testing.T, name string, fn func(t testing.TB)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Helper()
		tester := &testTester{T: t}
		fn(tester)
		if tester.failed == "" {
			t.Fatal("Should have failed")
		} else {
			t.Log(tester.failed)
		}
	})
}

func assertOk(t *testing.T, name string, fn func(t testing.TB)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Helper()
		tester := &testTester{T: t}
		fn(tester)
		if tester.failed != "" {
			t.Fatal("Should not have failed with:\n", tester.failed)
		}
	})
}
