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
	assertOk(t, func(t testing.TB) {
		Equal(t, Data{"expected", 1234}, Data{"expected", 1234})
	})
	assertFail(t, func(t testing.TB) {
		Equal(t, Data{"expected\ntext", 1234}, Data{"actual\ntext", 1234})
	})
}

func TestNotEqual(t *testing.T) {
	assertOk(t, func(t testing.TB) {
		NotEqual(t, Data{"expected", 1234}, Data{"expected", 1235})
	})
	assertFail(t, func(t testing.TB) {
		NotEqual(t, Data{"expected", 1234}, Data{"expected", 1234})
	})
}

func TestContains(t *testing.T) {
	assertOk(t, func(t testing.TB) {
		Contains(t, "a haystack with a needle in it", "needle")
	})
	assertFail(t, func(t testing.TB) {
		Contains(t, "a haystack with a needle in it", "screw")
	})
}

func TestNotContains(t *testing.T) {
	assertOk(t, func(t testing.TB) {
		NotContains(t, "a haystack with a needle in it", "screw")
	})
	assertFail(t, func(t testing.TB) {
		NotContains(t, "a haystack with a needle in it", "needle")
	})
}

func TestEqualError(t *testing.T) {
	assertOk(t, func(t testing.TB) {
		EqualError(t, fmt.Errorf("hello"), "hello")
	})
	assertFail(t, func(t testing.TB) {
		EqualError(t, fmt.Errorf("hello"), "goodbye")
	})
}

func TestZero(t *testing.T) {
	assertOk(t, func(t testing.TB) {
		Zero(t, Data{})
	})
	assertFail(t, func(t testing.TB) {
		Zero(t, Data{Str: "str"})
	})
}

func TestNotZero(t *testing.T) {
	assertOk(t, func(t testing.TB) {
		notZero := Data{Str: "hello"}
		NotZero(t, notZero)
	})
	assertFail(t, func(t testing.TB) {
		zero := Data{}
		NotZero(t, zero)
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

func assertFail(t *testing.T, fn func(testing.TB)) {
	t.Helper()
	tester := &testTester{T: t}
	fn(tester)
	if tester.failed == "" {
		t.Fatal("Should have failed with:\n" + tester.failed)
	} else {
		t.Log(tester.failed)
	}
}

func assertOk(t *testing.T, fn func(testing.TB)) {
	t.Helper()
	tester := &testTester{T: t}
	fn(tester)
	if tester.failed != "" {
		t.Fatal("Should not have failed with:\n", tester.failed)
	}
}
