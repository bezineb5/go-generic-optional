package opt

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Optional[T any] struct {
	value T
	exists bool
}

// New creates a new Optional without a value.
func New[T any]() Optional[T] {
	var empty T
	return Optional[T]{empty, false}
}

// Of creates a new Optional with a value.
func Of[T any](value T) Optional[T] {
	return Optional[T]{value, true}
}

// Get returns the value and whether it exists.
// It's invalid to use the returned value if the bool is false.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.exists
}

// GetOrElse returns the value if it exists and returns defaultValue otherwise.
func (o Optional[Value]) GetOrElse(defaultValue Value) Value {
	if !o.exists {
		return defaultValue
	}
	return o.value
}

// MustGet returns the value if it exists and panics otherwise.
func (o Optional[Value]) MustGet() Value {
	if !o.exists {
		panic(".MustGet() called on optional Optional value that doesn't exist.")
	}
	return o.value
}

// If allows you to handle an optional if it exists, otherwise return
func If[T any, R any](optional Optional[T], handler func(T) R) Optional[R] {
	if item, ok := optional.Get(); ok {
		return Of(handler(item))
	}
	return New[R]()
}

func (o Optional[Value]) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&struct {
			Value  Value `json:"value"`
			Exists bool  `json:"exists"`
		}{
			Value:  o.value,
			Exists: o.exists,
		},
	)
}

func (o *Optional[Value]) UnmarshalJSON(data []byte) error {
	s := &struct {
		Value  Value `json:"value"`
		Exists bool  `json:"exists"`
	}{}

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	o.value = s.Value
	o.exists = s.Exists

	return nil
}

// Scan implements the Scanner interface.
func (o *Optional[Value]) Scan(value any) error {
	if value == nil {
		o.exists = false
		return nil
	}

	o.exists = true
	var ok bool
	o.value, ok = value.(Value)
	if !ok {
		return fmt.Errorf("failed to scan a '%v' into an Optional", value)
	}

	return nil
}

// Value implements the Valuer interface.
func (o Optional[Value]) Value() (driver.Value, error) {
	if !o.exists {
		return nil, nil
	}
	return o.value, nil
}
