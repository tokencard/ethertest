package peek

import (
	"errors"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

func Call(path string, value interface{}, method string, args ...interface{}) ([]interface{}, error) {
	val, err := Peek(path, value)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(val)

	m := v.MethodByName(method)

	ar := []reflect.Value{}

	for _, av := range args {
		ar = append(ar, reflect.ValueOf(av))
	}

	r := m.Call(ar)

	res := []interface{}{}

	for _, rv := range r {
		res = append(res, rv.Interface())
	}

	return res, nil

}

func Peek(path string, value interface{}) (interface{}, error) {
	v := reflect.ValueOf(value)
	p := strings.Split(path, ".")
	_, r, err := peek(p, v)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func peek(path []string, v reflect.Value) (reflect.Value, interface{}, error) {

	if len(path) == 0 {
		return v, v.Interface(), nil
	}

	if len(path) == 1 {

		if path[0] == "" {
			return v, v.Interface(), nil
		}

		switch v.Kind() {
		case reflect.Ptr:
			return peek(path, v.Elem())
		case reflect.Interface:
			return peek(path, v.Elem())
		case reflect.Struct:
			fieldName := path[0]
			firstLetter, _ := utf8.DecodeRuneInString(fieldName)
			if unicode.IsLower(firstLetter) {
				cp := reflect.New(v.Type()).Elem()
				cp.Set(v)
				f := cp.FieldByName(fieldName)
				rf := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
				return rf, rf.Interface(), nil
			}
			f := v.FieldByName(fieldName)
			return f, f.Interface(), nil
		default:
			return reflect.ValueOf(nil), nil, errors.New("not yet implemeted")
		}
	}

	switch v.Kind() {
	case reflect.Ptr:
		return peek(path, v.Elem())
	case reflect.Interface:
		return peek(path, v.Elem())
	case reflect.Struct:
		fieldName := path[0]
		firstLetter, _ := utf8.DecodeRuneInString(fieldName)
		if unicode.IsLower(firstLetter) {
			cp := reflect.New(v.Type()).Elem()
			cp.Set(v)
			f := cp.FieldByName(fieldName)
			rf := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
			return peek(path[1:], rf)
		}
		f := v.FieldByName(fieldName)
		return peek(path[1:], f)
	default:
		return reflect.ValueOf(nil), nil, errors.New("not yet implemeted")
	}

}
