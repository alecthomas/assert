package example

import (
	"github.com/alecthomas/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

type Person struct {
	Name string
	Age  int
}

func TestDiff(t *testing.T) {
	expected := []*Person{{"Alec", 20}, {"Bob", 21}, {"Sally", 22}}
	actual := []*Person{{"Alex", 20}, {"Bob", 22}, {"Sally", 22}}
	assert.Equal(t, expected, actual)
}

func TestStretchrDiff(t *testing.T) {
	expected := []*Person{{"Alec", 20}, {"Bob", 21}, {"Sally", 22}}
	actual := []*Person{{"Alex", 20}, {"Bob", 22}, {"Sally", 22}}
	require.Equal(t, expected, actual)
}
