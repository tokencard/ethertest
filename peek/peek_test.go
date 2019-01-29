package peek_test

import (
	"testing"

	"github.com/tokencard/ethertest/peek"
	"github.com/stretchr/testify/assert"
)

func TestStructures(t *testing.T) {

	cases := []struct {
		name          string
		path          string
		data          interface{}
		expectedValue interface{}
		expectedError error
	}{
		{
			name:          "simple_struct_first_level_public",
			path:          "Foo",
			data:          struct{ Foo string }{Foo: "bar"},
			expectedValue: "bar",
			expectedError: nil,
		},
		{
			name:          "simple_struct_pointer_first_level_public",
			path:          "Foo",
			data:          &struct{ Foo string }{Foo: "bar"},
			expectedValue: "bar",
			expectedError: nil,
		},
		{
			name:          "simple_struct_first_level_private",
			path:          "foo",
			data:          struct{ foo string }{foo: "bar"},
			expectedValue: "bar",
			expectedError: nil,
		},
		{
			name:          "simple_struct_pointer_first_level_private",
			path:          "foo",
			data:          &struct{ foo string }{foo: "bar"},
			expectedValue: "bar",
			expectedError: nil,
		},
		{
			name: "struct_second_level_private",
			path: "foo.Bar",
			data: struct{ foo struct{ Bar string } }{
				foo: struct{ Bar string }{Bar: "baz"},
			},
			expectedValue: "baz",
			expectedError: nil,
		},
		{
			name: "struct_pointer_second_level_private",
			path: "foo.Bar",
			data: struct{ foo *struct{ Bar string } }{
				foo: &struct{ Bar string }{Bar: "baz"},
			},
			expectedValue: "baz",
			expectedError: nil,
		},
		{
			name: "struct_second_level",
			path: "Foo.Bar",
			data: struct{ Foo struct{ Bar string } }{
				Foo: struct{ Bar string }{Bar: "baz"},
			},
			expectedValue: "baz",
			expectedError: nil,
		},
		{
			name: "struct_pointer_second_level",
			path: "Foo.Bar",
			data: struct{ Foo *struct{ Bar string } }{
				Foo: &struct{ Bar string }{Bar: "baz"},
			},
			expectedValue: "baz",
			expectedError: nil,
		},
		{
			name: "struct_behind_interface_pointer_second_level",
			path: "Foo.Bar",
			data: struct{ Foo interface{} }{
				Foo: struct{ Bar string }{Bar: "baz"},
			},
			expectedValue: "baz",
			expectedError: nil,
		},
		{
			name: "struct_pointer_behind_interface_pointer_second_level",
			path: "Foo.Bar",
			data: struct{ Foo interface{} }{
				Foo: &struct{ Bar string }{Bar: "baz"},
			},
			expectedValue: "baz",
			expectedError: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert := assert.New(t)
			val, err := peek.Peek(c.path, c.data)
			assert.Equal(c.expectedError, err, "Unexpected error.")
			assert.Equal(c.expectedValue, val, "Unexpected value.")
		})
	}
}
