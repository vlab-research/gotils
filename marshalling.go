package gotils

import (
	"encoding/json"
	"reflect"
)

type CastMap map[string]func(interface{}) (interface{}, error)

func MarshalWithCasts(b []byte, emptyObj interface{}, castFns CastMap) (interface{}, error) {

	ttype := reflect.TypeOf(emptyObj)

	fields := []reflect.StructField{}
	for i := 0; i < ttype.NumField(); i++ {
		f := ttype.Field(i)
		if _, ok := castFns[f.Name]; ok {
			f.Type = reflect.TypeOf((*interface{})(nil)).Elem()
		}
		fields = append(fields, f)
	}
	t := reflect.StructOf(fields)

	obj := reflect.New(t).Interface()
	newObj := reflect.New(ttype).Elem()

	err := json.Unmarshal(b, obj)
	if err != nil {
		return newObj.Interface(), err
	}

	for i := 0; i < ttype.NumField(); i++ {
		ft := t.Field(i)
		fv := reflect.ValueOf(obj).Elem().Field(i)
		if fn, ok := castFns[ft.Name]; ok {
			s, err := fn(fv.Interface())
			if err != nil {
				return emptyObj, err
			}
			fv = reflect.ValueOf(s)
		}

		newObj.Field(i).Set(fv)
	}

	return newObj.Interface(), nil
}
