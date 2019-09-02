package uid

import (
	"fmt"
	"testing"
)

func TestUID(t *testing.T){
	 uid :=NewV1()
	 //t.Fatal(uid)
	 fmt.Println(uid)
	udata,err := uid.MarshalBinary()
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println("MarshalBinary:",udata)
	fmt.Println("uuid version:",uid.Version())
}
