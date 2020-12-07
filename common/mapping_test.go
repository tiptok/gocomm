package common

import "testing"

func Test_ObjectToMap(t *testing.T) {
	var o = struct {
		A  int
		A1 int
		B  string
		B1 string
	}{
		A1: 5,
		B1: "10",
	}

	//type O1 struct {
	//	A int
	//	A1 int
	//	B string
	//	B1 string
	//}
	//var o =&O1{
	//	A1:5,
	//	B1:"10",
	//}
	for k, v := range ObjectToMap(o) {
		t.Log(k, v)
	}
}

func TestLoadCustomField(t *testing.T) {

	type User struct {
		Name string
		Id   int
	}
	v := []struct {
		Name string
		Id   int
	}{{Name: "c1", Id: 1}, {Name: "c2", Id: 2}, {Name: "c3", Id: 3}}
	ret := LoadCustomField(&v, "Name")
	t.Log(JsonAssertString(ret))

	v2 := struct {
		Name string
		Id   int
	}{Name: "c1", Id: 1}
	ret2 := LoadCustomField(v2, "Name", "Id")
	t.Log(JsonAssertString(ret2))

	v3 := []*User{&User{Name: "c1", Id: 1}, &User{Name: "c2", Id: 2}, &User{Name: "c3", Id: 3}}
	ret3 := LoadCustomField(&v3, "Name", "Name2")
	t.Log(JsonAssertString(ret3))
}

func TestAppendCustomField(t *testing.T) {
	customMap := map[string]interface{}{"1": "h", "2": "e"}
	t.Log(customMap)

	customeFiles := map[string]interface{}{
		"Value":  map[string]interface{}{"a": "a", "b": "b"},
		"Value2": "cc",
		"Value3": 9999,
	}
	t.Log(AppendCustomField(customMap, customeFiles))

	t.Log(AppendCustomField(struct {
		Name string `json:"name"`
		Age  int64  `json:"age"`
	}{Name: "ccc", Age: 20}, customeFiles))
}
