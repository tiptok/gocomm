package xstring

import (
	"fmt"
	"github.com/tiptok/gocomm/common"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	//letterIdxBits  = 6 // 6 bits to represent a letter index
	idLen = 8
	//defaultRandLen = 8
	//letterIdxMask  = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	//letterIdxMax   = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandId() string {
	b := []byte(common.RandomStringWithChars(idLen, letterBytes))
	return fmt.Sprintf("%x%x%x%x", b[0:2], b[2:4], b[4:6], b[6:8])
}

func TakeWithPriority(fns ...func() string) string {
	for _, fn := range fns {
		val := fn()
		if len(val) > 0 {
			return val
		}
	}

	return ""
}
