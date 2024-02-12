package opt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSmoke(t *testing.T) {
	require.True(t, true)
}

func TestNew(t *testing.T) {
	anOptional := New[string]()
	_, ok := anOptional.Get()
	require.Equal(t, ok, false)
	require.False(t, anOptional.Exists())
	require.Equal(t, "something", anOptional.GetOrElse("something"))
}

func TestOf(t *testing.T) {
	anOptional := Of("hello")
	require.True(t, anOptional.Exists())
	value, ok := anOptional.Get()
	require.Equal(t, true, ok)
	require.Equal(t, "hello", value)
	require.Equal(t, "hello", anOptional.MustGet())
}

func TestMustGet(t *testing.T) {
	anOptional := New[string]()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	anOptional.MustGet()
}

func TestIf(t *testing.T) {
	anOptional := New[string]()

	ranIf := false

	If(anOptional, func(value string) bool {
		ranIf = true
		return true
	})

	require.False(t, ranIf)

	anOptional = Of("hello")

	If(anOptional, func(value string) bool {
		ranIf = true
		return true
	})

	require.True(t, ranIf)
}

func TestMarshall(t *testing.T) {
	anOptional := New[string]()

	data, err := anOptional.MarshalJSON()

	require.Nil(t, err)
	require.Equal(t, `null`, string(data))

	anOptional = Of("hello")

	data, err = anOptional.MarshalJSON()

	require.Nil(t, err)
	require.Equal(t, `"hello"`, string(data))
}

func TestUnmarshall(t *testing.T) {
	anOptional := New[string]()

	err := anOptional.UnmarshalJSON([]byte(`null`))
	require.Nil(t, err)
	require.False(t, anOptional.Exists())

	anOptional = New[string]()
	err = anOptional.UnmarshalJSON([]byte(`"hello"`))
	require.Nil(t, err)
	require.True(t, anOptional.Exists())
	require.Equal(t, "hello", anOptional.MustGet())

	err = anOptional.UnmarshalJSON([]byte(`"asdjl:1k2l;j'""`))
	require.NotNil(t, err)
}

func TestFlatMap(t *testing.T) {
	// Test when the optional is empty
	emptyOptional := New[string]()
	mappedOptional := FlatMap(emptyOptional, func(item string) Optional[int] {
		return Of(len(item))
	})
	require.False(t, mappedOptional.Exists())

	// Test when the optional is not empty
	nonEmptyOptional := Of("hello")
	mappedOptional = FlatMap(nonEmptyOptional, func(item string) Optional[int] {
		return Of(len(item))
	})
	require.True(t, mappedOptional.Exists())
	require.Equal(t, 5, mappedOptional.MustGet())

	// Test when the optional is not empty and the mapper returns an empty optional
	mappedOptional = FlatMap(nonEmptyOptional, func(item string) Optional[int] {
		return New[int]()
	})
	require.False(t, mappedOptional.Exists())
}

func TestOrElse(t *testing.T) {
	// Test when the optional is empty
	emptyOptional := New[string]()
	defaultValue := "default"
	result := emptyOptional.OrElse(defaultValue)
	require.True(t, result.Exists())
	require.Equal(t, defaultValue, result.MustGet())

	// Test when the optional is not empty
	nonEmptyOptional := Of("hello")
	defaultValue = "default"
	result = nonEmptyOptional.OrElse(defaultValue)
	require.True(t, result.Exists())
	require.Equal(t, "hello", result.MustGet())
}

func TestFilter(t *testing.T) {
	// Test when the optional is empty
	emptyOptional := New[string]()
	filteredOptional := emptyOptional.Filter(func(item string) bool {
		return len(item) > 0
	})
	require.False(t, filteredOptional.Exists())

	// Test when the optional is not empty and the predicate returns true
	nonEmptyOptional := Of("hello")
	filteredOptional = nonEmptyOptional.Filter(func(item string) bool {
		return len(item) > 0
	})
	require.True(t, filteredOptional.Exists())
	require.Equal(t, "hello", filteredOptional.MustGet())

	// Test when the optional is not empty and the predicate returns false
	filteredOptional = nonEmptyOptional.Filter(func(item string) bool {
		return len(item) > 10
	})
	require.False(t, filteredOptional.Exists())
}
