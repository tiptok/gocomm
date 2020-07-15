package common

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	str := "123456"
	md := MD5String([]byte(str))
	//t.Log("str", str)
	//t.Log("str MD5 ", md)
	ok := MD5Verify(str, md)
	if !ok {
		t.Fatal(fmt.Sprintf("MD5Verify err . in : %v except: %v out:%v", str, true, md))
	}
}

func TestAes(t *testing.T) {
	key := []byte("0123456789abcdef") //key len  must 16/24/32
	result, err := AesEncrypt([]byte("hello world"), key)
	if err != nil {
		panic(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(result))
	origData, err := AesDecrypt(result, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(origData))
}
