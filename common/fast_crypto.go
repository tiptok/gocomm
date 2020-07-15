package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"strings"
)

//md5 加密/校验
//md5加密
func MD5String(code []byte) string {
	h := md5.New()
	h.Write(code)
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

//md5校验
func MD5Verify(code string, md5Str string) bool {
	codeMD5 := MD5String([]byte(code))
	return 0 == strings.Compare(codeMD5, md5Str)
}

// AES 加密/解密
func AesEncrypt(plantText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plantText = pKCS7Padding(plantText, block.BlockSize())
	//偏转向量iv长度等于密钥key块大小
	iv := key[:block.BlockSize()]
	blockModel := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(plantText))

	blockModel.CryptBlocks(cipherText, plantText)
	return cipherText, nil
}
func AesDecrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//偏转向量iv长度等于密钥key块大小
	iv := key[:block.BlockSize()]
	blockModel := cipher.NewCBCDecrypter(block, iv)
	plantText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plantText, cipherText)
	plantText = pKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
}
func pKCS7UnPadding(plantText []byte, blockSize int) []byte {
	//AES Decrypt pkcs7padding CBC, key for choose algorithm
	length := len(plantText)
	unPadding := int(plantText[length-1])
	return plantText[:(length - unPadding)]
}
func pKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}
