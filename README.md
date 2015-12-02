# Go assertion library (fork of [stretchr/testify/assert](https://github.com/stretchr/testify/tree/master/assert))

This is a fork of stretchr/testify/assert that fails the current test function when an assertion fails. In my
opinion this is more aligned with the concept of an assertion, making the package more intuitive. It
also makes for much cleaner code, as it does not require an if statement for every potentially
failing assertion. eg.

Idiomatic usage of stretchr/testify/assert:

```go
v, err := pkg.Func()
if assert.NoError(t, err) {
  if assert.NotNil(t, v) {
    if assert.Equal(t, v.Value, 1) {
      // ...
    }
  }
}
```

Idiomatic usage of this package:

```go
v, err := pkg.Func()
assert.NoError(t, err)
assert.NotNil(t, v)
assert.Equal(t, v.Value, 1)
```
