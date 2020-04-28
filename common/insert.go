package common

import "reflect"

func Insert(slice interface{}, pos int, value interface{}) interface{}  {
	v := reflect.ValueOf(slice)
	insert := reflect.ValueOf(value)
	if insert.Len() == 0 {
		return v.Interface()
	}
	ne := reflect.MakeSlice(reflect.TypeOf(value), insert.Len(), insert.Len())
	for i := 0; i < insert.Len(); i++ {
		ne.Index(i).Set(insert.Index(i))
	}
	v = reflect.AppendSlice(v.Slice(0, pos), reflect.AppendSlice(ne, v.Slice(pos, v.Len())))
	return v.Interface()
}
