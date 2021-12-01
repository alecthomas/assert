// Package assert provides type-safe assertions with clean error messages.
package assert

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/google/go-cmp/cmp"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
)

// Compare two values for equality and return true or false.
func Compare[T any](t testing.TB, x, y T) bool {
	return cmp.Equal(x, y)
}

// Equal asserts that "expected" and "actual" are equal.
//
// If they are not a diff of the Go representation of the values will be displayed.
func Equal[T any](t testing.TB, expected, actual T, msgAndArgs ...interface{}) {
	if cmp.Equal(expected, actual) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected values to be equal:", msgAndArgs...)
	t.Fatalf("%s\n%s", msg, diff(expected, actual))
}

// NotEqual asserts that "expected" is not equal to "actual".
//
// If they are not a diff of the Go representation of the values will be displayed.
func NotEqual[T any](t testing.TB, expected, actual T, msgAndArgs ...interface{}) {
	if !cmp.Equal(expected, actual) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected values to not be equal:", msgAndArgs...)
	t.Fatalf("%s\n%s", msg, repr.String(expected, repr.Indent("  ")))
}

// Contains asserts that "slice" contains "element".
func Contains[T any](t testing.TB, slice []T, element T, msgAndArgs ...interface{}) {
	for _, el := range slice {
		if cmp.Equal(el, element) {
			return
		}
	}
	t.Helper()
	var msg string
	if len(msgAndArgs) == 0 {
		msg = fmt.Sprintf("%s does not contain %s", repr.String(slice), repr.String(element))
	} else {
		msg = formatMsgAndArgs("", msgAndArgs...)
	}
	t.Fatal(msg)
}

// NotContains asserts that "slice" does not contain "element".
func NotContains[T any](t testing.TB, slice []T, element T, msgAndArgs ...interface{}) {
	found := false
	for _, el := range slice {
		if cmp.Equal(el, element) {
			found = true
			break
		}
	}
	if !found {
		return
	}
	t.Helper()
	var msg string
	if len(msgAndArgs) == 0 {
		msg = fmt.Sprintf("%s should not contain %s", repr.String(slice), repr.String(element))
	} else {
		msg = formatMsgAndArgs("", msgAndArgs...)
	}
	t.Fatal(msg)
}

// ContainsKey asserts that a map contains the given key.
func ContainsKey[K comparable, V any](t testing.TB, amap map[K]V, key K, msgAndArgs ...interface{}) {
	found := false
	for k := range amap {
		if cmp.Equal(k, key) {
			found = true
			break
		}
	}
	if found {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected key to be present in map:", msgAndArgs...)
	t.Fatalf("%s\nKey: %s\nMap: %s\n", msg, repr.String(key), repr.String(amap))
}

// ContainsValue asserts that a map contains the given value.
func ContainsValue[K comparable, V any](t testing.TB, amap map[K]V, value V, msgAndArgs ...interface{}) {
	found := false
	for _, v := range amap {
		if cmp.Equal(v, value) {
			found = true
			break
		}
	}
	if found {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected value to be present in map:", msgAndArgs...)
	t.Fatalf("%s\nValue: %s\nMap: %s\n", msg, repr.String(value), repr.String(amap))
}

// AllMap asserts that all entries in a map pass the given predicate.
func AllMap[K comparable, V any](t testing.TB, amap map[K]V, predicate func(k K, v V) bool, msgAndArgs ...interface{}) {
	for k, v := range amap {
		if !predicate(k, v) {
			t.Helper()
			msg := formatMsgAndArgs("Not all values in map were true:", msgAndArgs...)
			t.Fatalf("%s\n%s", msg, repr.String(amap))
		}
	}
}

// AllMap asserts that at least one entry in a map pass the given predicate.
func AnyMap[K comparable, V any](t testing.TB, amap map[K]V, predicate func(k K, v V) bool, msgAndArgs ...interface{}) {
	for k, v := range amap {
		if predicate(k, v) {
			return
		}
	}
	t.Helper()
	msg := formatMsgAndArgs("", msgAndArgs...)
	t.Fatalf("%s\n%s", msg, repr.String(amap))
}

// AnySlice asserts that predicate is true for at least one element of slice.
func AnySlice[T any](t testing.TB, slice []T, predicate func(i int, el T) bool, msgAndArgs ...interface{}) {
	for i, el := range slice {
		if predicate(i, el) {
			return
		}
	}
	t.Helper()
	msg := formatMsgAndArgs("No elements in slice matched:", msgAndArgs...)
	t.Fatalf("%s\n%s", msg, repr.String(slice))
}

// AllSlice asserts that predicate is true for all elements of slice.
func AllSlice[T any](t testing.TB, slice []T, predicate func(i int, el T) bool, msgAndArgs ...interface{}) {
	for i, el := range slice {
		if !predicate(i, el) {
			t.Helper()
			msg := formatMsgAndArgs("Not all elements in slice matched:", msgAndArgs...)
			t.Fatalf("%s\n%s", msg, repr.String(slice))
		}
	}
}

// Zero asserts that a value is its zero value.
func Zero[T any](t testing.TB, value T, msgAndArgs ...interface{}) {
	var zero T
	if cmp.Equal(value, zero) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected a zero value:", msgAndArgs...)
	t.Fatalf("%s\n%s", msg, repr.String(value, repr.Indent("  ")))
}

// NotZero asserts that a value is not its zero value.
func NotZero[T any](t testing.TB, value T, msgAndArgs ...interface{}) {
	var zero T
	if !cmp.Equal(value, zero) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Did not expect a zero value:", msgAndArgs...)
	t.Fatalf("%s\n%s", msg, repr.String(value, repr.Indent("  ")))
}

type floats interface{ float32 | float64 }

// AlmostEqual asserts that two floats are almost equal, within delta.
func AlmostEqual[T floats](t testing.TB, lhs, rhs T, delta T, msgAndArgs ...interface{}) {
	if math.Abs(float64(lhs-rhs)) <= float64(delta) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected values to be almost equal:", msgAndArgs...)
	t.Fatalf("%s\n%f ~= %f", msg, float64(lhs), float64(rhs))
}

// NotAlmostEqual asserts that two floats are almost equal, within delta.
func NotAlmostEqual[T floats](t testing.TB, lhs, rhs T, delta T, msgAndArgs ...interface{}) {
	if math.Abs(float64(lhs-rhs)) >= float64(delta) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected values to not be almost equal:", msgAndArgs...)
	t.Fatalf("%s\n%f ~= %f", msg, float64(lhs), float64(rhs))
}

// NoError asserts that an error is nil.
func Error(t testing.TB, err error, msgAndArgs ...interface{}) {
	if err != nil {
		return
	}
	t.Helper()
	t.Fatal(formatMsgAndArgs("Expected an error", msgAndArgs...))
}

// NoError asserts that an error is nil.
func NoError(t testing.TB, err error, msgAndArgs ...interface{}) {
	if err == nil {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Did not expect an error but got:", msgAndArgs...)
	t.Fatalf("%s\n%s", msg, err)
}

// True asserts that an expression is true.
func True(t testing.TB, ok bool, msgAndArgs ...interface{}) {
	if ok {
		return
	}
	t.Helper()
	t.Fatal("Expected expression to be true")
}

// False asserts that an expression is false.
func False(t testing.TB, ok bool, msgAndArgs ...interface{}) {
	if !ok {
		return
	}
	t.Helper()
	t.Fatal("Expected expression to be false")
}

func diff[T any](lhs, rhs T) string {
	lhss := repr.String(lhs, repr.Indent("  ")) + "\n"
	rhss := repr.String(rhs, repr.Indent("  ")) + "\n"
	edits := myers.ComputeEdits("a.txt", lhss, rhss)
	lines := strings.Split(fmt.Sprint(gotextdiff.ToUnified("expected.txt", "actual.txt", lhss, edits)), "\n")
	return strings.Join(lines[3:], "\n")
}

func formatMsgAndArgs(dflt string, msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 {
		return dflt
	}
	return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
}
