// Package assert provides type-safe assertions with clean error messages.
package assert

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
)

// A CompareOption modifies how object comparisons behave.
type CompareOption func() []repr.Option

// Exclude fields of the given type from comparison.
func Exclude[T any]() CompareOption {
	return func() []repr.Option {
		return []repr.Option{repr.Hide[T]()}
	}
}

// OmitEmpty fields from comparison.
func OmitEmpty() CompareOption {
	return func() []repr.Option {
		return []repr.Option{repr.OmitEmpty(true)}
	}
}

// IgnoreGoStringer ignores GoStringer implementations when comparing.
func IgnoreGoStringer() CompareOption {
	return func() []repr.Option {
		return []repr.Option{repr.IgnoreGoStringer()}
	}
}

// Compare two values for equality and return true or false.
func Compare[T any](t testing.TB, x, y T, options ...CompareOption) bool {
	return objectsAreEqual(x, y, options...)
}

func extractCompareOptions(msgAndArgs ...any) ([]any, []CompareOption) {
	compareOptions := []CompareOption{}
	out := []any{}
	for _, arg := range msgAndArgs {
		if opt, ok := arg.(CompareOption); ok {
			compareOptions = append(compareOptions, opt)
		} else {
			out = append(out, arg)
		}
	}
	return out, compareOptions
}

// HasPrefix asserts that the string s starts with prefix.
func HasPrefix(t testing.TB, s, prefix string, msgAndArgs ...any) {
	if strings.HasPrefix(s, prefix) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected string to have prefix:", msgAndArgs...)
	t.Fatalf("%s\nPrefix: %q\nString: %q\n", msg, prefix, s)
}

// HasSuffix asserts that the string s ends with suffix.
func HasSuffix(t testing.TB, s, suffix string, msgAndArgs ...any) {
	if strings.HasSuffix(s, suffix) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected string to have suffix:", msgAndArgs...)
	t.Fatalf("%s\nSuffix: %q\nString: %q\n", msg, suffix, s)
}

// Equal asserts that "expected" and "actual" are equal.
//
// If they are not, a diff of the Go representation of the values will be displayed.
func Equal[T any](t testing.TB, expected, actual T, msgArgsAndCompareOptions ...any) {
	msgArgsAndCompareOptions, compareOptions := extractCompareOptions(msgArgsAndCompareOptions...)
	if objectsAreEqual(expected, actual, compareOptions...) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected values to be equal:", msgArgsAndCompareOptions...)
	t.Fatalf("%s\n%s", msg, Diff(expected, actual, compareOptions...))
}

// NotEqual asserts that "expected" is not equal to "actual".
//
// If they are equal the expected value will be displayed.
func NotEqual[T any](t testing.TB, expected, actual T, msgArgsAndCompareOptions ...any) {
	msgArgsAndCompareOptions, compareOptions := extractCompareOptions(msgArgsAndCompareOptions...)
	if !objectsAreEqual(expected, actual, compareOptions...) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected values to not be equal but both were:", msgArgsAndCompareOptions...)
	t.Fatalf("%s\n%s", msg, repr.String(expected, repr.Indent("  ")))
}

// Contains asserts that "haystack" contains "needle".
func Contains(t testing.TB, haystack string, needle string, msgAndArgs ...any) {
	if strings.Contains(haystack, needle) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Haystack does not contain needle.", msgAndArgs...)
	t.Fatalf("%s\nNeedle: %q\nHaystack: %q\n", msg, needle, haystack)
}

// NotContains asserts that "haystack" does not contain "needle".
func NotContains(t testing.TB, haystack string, needle string, msgAndArgs ...any) {
	if !strings.Contains(haystack, needle) {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Haystack should not contain needle.", msgAndArgs...)
	quotedHaystack, quotedNeedle, positions := needlePosition(haystack, needle)
	t.Fatalf("%s\nNeedle: %s\nHaystack: %s\n          %s\n", msg, quotedNeedle, quotedHaystack, positions)
}

// SliceContains asserts that "haystack" contains "needle".
func SliceContains[T any](t testing.TB, haystack []T, needle T, msgAndArgs ...interface{}) {
	t.Helper()
	for _, item := range haystack {
		if objectsAreEqual(item, needle) {
			return
		}
	}

	msg := formatMsgAndArgs("Haystack does not contain needle.", msgAndArgs...)
	needleRepr := repr.String(needle, repr.Indent("  "))
	haystackRepr := repr.String(haystack, repr.Indent("  "))
	t.Fatalf("%s\nNeedle: %s\nHaystack: %s\n", msg, needleRepr, haystackRepr)
}

// NotSliceContains asserts that "haystack" does not contain "needle".
func NotSliceContains[T any](t testing.TB, haystack []T, needle T, msgAndArgs ...interface{}) {
	t.Helper()
	for _, item := range haystack {
		if objectsAreEqual(item, needle) {
			msg := formatMsgAndArgs("Haystack should not contain needle.", msgAndArgs...)
			needleRepr := repr.String(needle, repr.Indent("  "))
			haystackRepr := repr.String(haystack, repr.Indent("  "))
			t.Fatalf("%s\nNeedle: %s\nHaystack: %s\n", msg, needleRepr, haystackRepr)
		}
	}
}

// Zero asserts that a value is its zero value.
func Zero[T any](t testing.TB, value T, msgAndArgs ...any) {
	var zero T
	if objectsAreEqual(value, zero) {
		return
	}
	val := reflect.ValueOf(value)
	if (val.Kind() == reflect.Slice || val.Kind() == reflect.Map || val.Kind() == reflect.Array) && val.Len() == 0 {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Expected a zero value but got:", msgAndArgs...)
	t.Fatalf("%s\n%s", msg, repr.String(value, repr.Indent("  ")))
}

// NotZero asserts that a value is not its zero value.
func NotZero[T any](t testing.TB, value T, msgAndArgs ...any) {
	var zero T
	if !objectsAreEqual(value, zero) {
		val := reflect.ValueOf(value)
		if !((val.Kind() == reflect.Slice || val.Kind() == reflect.Map || val.Kind() == reflect.Array) && val.Len() == 0) {
			return
		}
	}
	t.Helper()
	msg := formatMsgAndArgs("Did not expect the zero value:", msgAndArgs...)
	t.Fatalf("%s\n%s", msg, repr.String(value))
}

// EqualError asserts that either an error is non-nil and that its message is what is expected,
// or that error is nil if the expected message is empty.
func EqualError(t testing.TB, err error, errString string, msgAndArgs ...any) {
	if err == nil && errString == "" {
		return
	}
	t.Helper()
	if err == nil {
		t.Fatal(formatMsgAndArgs("Expected an error", msgAndArgs...))
	}
	if err.Error() != errString {
		msg := formatMsgAndArgs("Error message not as expected:", msgAndArgs...)
		t.Fatalf("%s\n%s", msg, Diff(errString, err.Error()))
	}
}

// IsError asserts than any error in "err"'s tree matches "target".
func IsError(t testing.TB, err, target error, msgAndArgs ...any) {
	if errors.Is(err, target) {
		return
	}
	t.Helper()
	t.Fatal(formatMsgAndArgs(fmt.Sprintf("Error tree %+v should contain error %q", err, target), msgAndArgs...))
}

// NotIsError asserts than no error in "err"'s tree matches "target".
func NotIsError(t testing.TB, err, target error, msgAndArgs ...any) {
	if !errors.Is(err, target) {
		return
	}
	t.Helper()
	t.Fatal(formatMsgAndArgs(fmt.Sprintf("Error tree %+v should NOT contain error %q", err, target), msgAndArgs...))
}

// Error asserts that an error is not nil.
func Error(t testing.TB, err error, msgAndArgs ...any) {
	if err != nil {
		return
	}
	t.Helper()
	t.Fatal(formatMsgAndArgs("Expected an error", msgAndArgs...))
}

// NoError asserts that an error is nil.
func NoError(t testing.TB, err error, msgAndArgs ...any) {
	if err == nil {
		return
	}
	t.Helper()
	msg := formatMsgAndArgs("Did not expect an error but got:", msgAndArgs...)
	t.Fatalf("%s\n%+v", msg, err)
}

// True asserts that an expression is true.
func True(t testing.TB, ok bool, msgAndArgs ...any) {
	if ok {
		return
	}
	t.Helper()
	t.Fatal(formatMsgAndArgs("Expected expression to be true", msgAndArgs...))
}

// False asserts that an expression is false.
func False(t testing.TB, ok bool, msgAndArgs ...any) {
	if !ok {
		return
	}
	t.Helper()
	t.Fatal(formatMsgAndArgs("Expected expression to be false", msgAndArgs...))
}

// Panics asserts that the given function panics.
func Panics(t testing.TB, fn func(), msgAndArgs ...any) {
	t.Helper()
	defer func() {
		if recover() == nil {
			msg := formatMsgAndArgs("Expected function to panic", msgAndArgs...)
			t.Fatal(msg)
		}
	}()
	fn()
}

// NotPanics asserts that the given function does not panic.
func NotPanics(t testing.TB, fn func(), msgAndArgs ...any) {
	t.Helper()
	defer func() {
		if err := recover(); err != nil {
			msg := formatMsgAndArgs("Expected function not to panic", msgAndArgs...)
			t.Fatalf("%s\nPanic: %v", msg, err)
		}
	}()
	fn()
}

// Diff returns a unified diff of the string representation of two values.
func Diff[T any](before, after T, compareOptions ...CompareOption) string {
	var lhss, rhss string
	// Special case strings so we get nice diffs.
	if l, ok := any(before).(string); ok {
		lhss = l + "\n"
		rhss = any(after).(string) + "\n"
	} else {
		ropts := expandCompareOptions(compareOptions...)
		lhss = repr.String(before, ropts...) + "\n"
		rhss = repr.String(after, ropts...) + "\n"
	}
	edits := myers.ComputeEdits("a.txt", lhss, rhss)
	lines := strings.Split(fmt.Sprint(gotextdiff.ToUnified("expected.txt", "actual.txt", lhss, edits)), "\n")
	if len(lines) < 3 {
		return ""
	}
	return strings.Join(lines[3:], "\n")
}

func formatMsgAndArgs(dflt string, msgAndArgs ...any) string {
	if len(msgAndArgs) == 0 {
		return dflt
	}
	format, ok := msgAndArgs[0].(string)
	if !ok {
		panic("message argument to assert function must be a fmt string")
	}
	return fmt.Sprintf(format, msgAndArgs[1:]...)
}

func needlePosition(haystack, needle string) (quotedHaystack, quotedNeedle, positions string) {
	quotedNeedle = strconv.Quote(needle)
	quotedNeedle = quotedNeedle[1 : len(quotedNeedle)-1]
	quotedHaystack = strconv.Quote(haystack)
	rawPositions := strings.ReplaceAll(quotedHaystack, quotedNeedle, strings.Repeat("^", len(quotedNeedle)))
	for _, rn := range rawPositions {
		if rn != '^' {
			positions += " "
		} else {
			positions += "^"
		}
	}
	return
}

func expandCompareOptions(options ...CompareOption) []repr.Option {
	ropts := []repr.Option{repr.Indent("  ")}
	for _, option := range options {
		ropts = append(ropts, option()...)
	}
	return ropts
}

func objectsAreEqual(expected, actual any, options ...CompareOption) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}
	if exp, eok := expected.([]byte); eok {
		if act, aok := actual.([]byte); aok {
			return bytes.Equal(exp, act)
		}
	}
	if exp, eok := expected.(string); eok {
		if act, aok := actual.(string); aok {
			return exp == act
		}
	}

	ropts := expandCompareOptions(options...)
	expectedStr := repr.String(expected, ropts...)
	actualStr := repr.String(actual, ropts...)

	return expectedStr == actualStr
}
