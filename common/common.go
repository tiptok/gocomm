package common

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must2 panics if the second parameter is not nil, otherwise returns the first parameter.
func Must2(v interface{}, err error) interface{} {
	Must(err)
	return v
}

// Error2 returns the err from the 2nd parameter.
func Error2(v interface{}, err error) error {
	return err
}

func LogF(format string, args interface{}) string {
	return fmt.Sprintf(format, args)
}

var randomChars = "ABCDEFGHJKMNPQRSTWXYZabcdefhjkmnprstwxyz2345678" /****默认去掉了容易混淆的字符oOLl,9gq,Vv,Uu,I1****/
func RandomString(l int) string {
	return RandomStringWithChars(l, randomChars)
}

func RandomStringWithChars(l int, chars string) string {
	if l <= 0 {
		return ""
	}
	if len(chars) == 0 {
		return ""
	}
	lenChars := len(chars) - 1
	rsp := bytes.NewBuffer(nil)
	rand.Seed(time.Now().Unix())
	for i := 0; i < l; i++ {
		rsp.WriteByte(chars[rand.Intn(lenChars)])
	}
	return rsp.String()
}
