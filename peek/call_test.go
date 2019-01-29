package peek_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokencard/ethertest/peek"
)

type pointerCallee struct {
}

func (c *pointerCallee) CallMe(foo string, bar int) error {
	return nil
}

type directCallee struct {
}

func (c directCallee) CallMe(foo string, bar int) error {
	return nil
}

func TestCall(t *testing.T) {

	cases := []struct {
		name           string
		path           string
		data           interface{}
		method         string
		args           []interface{}
		expectedValues []interface{}
		expectedError  error
	}{
		{
			name:           "simple_struct_first_level_pointer_private",
			path:           "",
			data:           &pointerCallee{},
			method:         "CallMe",
			args:           []interface{}{"baz", 42},
			expectedValues: []interface{}{nil},
			expectedError:  nil,
		},
		{
			name:           "simple_struct_first_level_direct_private",
			path:           "",
			data:           &directCallee{},
			method:         "CallMe",
			args:           []interface{}{"baz", 42},
			expectedValues: []interface{}{nil},
			expectedError:  nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert := assert.New(t)
			vals, err := peek.Call(c.path, c.data, c.method, c.args...)
			assert.Equal(c.expectedError, err, "Unexpected error.")
			assert.Equal(c.expectedValues, vals, "Unexpected values.")
		})
	}
}
