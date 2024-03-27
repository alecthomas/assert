package assert

import (
	"fmt"
	"os"
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
	assertOk(t, "Zero length byte arrays", func(t testing.TB) {
		Equal(t, []byte(""), []byte(nil))
	})
	assertOk(t, "Identical byte arrays", func(t testing.TB) {
		Equal(t, []byte{4, 2}, []byte{4, 2})
	})
	assertOk(t, "Identical numbers", func(t testing.TB) {
		Equal(t, 42, 42)
	})
	assertFail(t, "DifferentStruct", func(t testing.TB) {
		Equal(t, Data{"expected\ntext", 1234}, Data{"actual\ntext", 1234})
	})
	assertFail(t, "Different bytes arrays", func(t testing.TB) {
		Equal(t, []byte{2, 4}, []byte{4, 2})
	})
	assertFail(t, "Different numbers", func(t testing.TB) {
		Equal(t, 42, 43)
	})
	assertOk(t, "Exclude", func(t testing.TB) {
		Equal(t, Data{Str: "expected", Num: 1234}, Data{Str: "expected"}, Exclude[int64]())
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
	assertFail(t, "Exclude", func(t testing.TB) {
		NotEqual(t, Data{Str: "expected", Num: 1234}, Data{Str: "expected"}, Exclude[int64]())
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

func TestError(t *testing.T) {
	assertOk(t, "Error", func(t testing.TB) {
		Error(t, fmt.Errorf("hello"))
	})
	assertFail(t, "Nil", func(t testing.TB) {
		Error(t, nil)
	})
}

func TestNoError(t *testing.T) {
	assertOk(t, "Nil", func(t testing.TB) {
		NoError(t, nil)
	})
	assertFail(t, "Error", func(t testing.TB) {
		NoError(t, fmt.Errorf("hello"))
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

func TestIsError(t *testing.T) {
	assertOk(t, "SameError", func(t testing.TB) {
		IsError(t, fmt.Errorf("os error: %w", os.ErrClosed), os.ErrClosed)
	})
	assertFail(t, "DifferentError", func(t testing.TB) {
		IsError(t, fmt.Errorf("not an os error"), os.ErrClosed)
	})
}

func TestInvalidFormatMsg(t *testing.T) {
	Panics(t, func() {
		NotZero(t, Data{}, 123)
	})
}

func TestNotIsError(t *testing.T) {
	assertFail(t, "SameError", func(t testing.TB) {
		NotIsError(t, fmt.Errorf("os error: %w", os.ErrClosed), os.ErrClosed)
	})
	assertOk(t, "DifferentError", func(t testing.TB) {
		NotIsError(t, fmt.Errorf("not an os error"), os.ErrClosed)
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
