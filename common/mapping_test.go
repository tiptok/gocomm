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
