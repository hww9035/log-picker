package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Md5ByString 字符串md5加密
func Md5ByString(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		panic(err)
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}

// Md5ByBytes 字节md5加密
func Md5ByBytes(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

// HashSha256StrToBytes 字符串形式hash256加密，返回字节数据
func HashSha256StrToBytes(input, key string) []byte {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(input))

	return m.Sum(nil)
}

// HashSha256BytesToBytes 字节形式hash256加密，返回字节数据
func HashSha256BytesToBytes(input []byte, key string) []byte {
	m := hmac.New(sha256.New, []byte(key))
	m.Write(input)

	return m.Sum(nil)
}

// HashSha256StrToHex 字符串形式hash256加密，返回16进制字符串
func HashSha256StrToHex(input, key string) string {
	return hex.EncodeToString(HashSha256StrToBytes(input, key))
}

// HashSha256BytesToHex 字节形式hash256加密，返回16进制字符串
func HashSha256BytesToHex(input []byte, key string) string {
	return hex.EncodeToString(HashSha256BytesToBytes(input, key))
}

// Base64StdEncodeBytes base64标准加密
func Base64StdEncodeBytes(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

// Base64UrlEncodeBytes base64URL加密
func Base64UrlEncodeBytes(input []byte) string {
	return base64.URLEncoding.EncodeToString(input)
}

// Base64StdDecode base64标准解密
func Base64StdDecode(input string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(input)
}

// Base64UrlDecode base64URL解密
func Base64UrlDecode(input string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(input)
}

// PKCS7Padding 填充字符串
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padText...)
}

// PKCS7UnPadding 删除填充字符串
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("origData error")
	}
	unPadding := int(origData[length-1])
	if length < unPadding {
		return []byte(""), nil
	}

	return origData[:(length - unPadding)], nil
}

// AesEncrypt AES加密
func AesEncrypt(input []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	// 对数据进行填充
	input = PKCS7Padding(input, blockSize)
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypt := make([]byte, len(input))
	blocMode.CryptBlocks(crypt, input)

	return crypt, nil
}

// AesDeCrypt AES解密
func AesDeCrypt(input []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(input))
	blockMode.CryptBlocks(origData, input)
	// 去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}

	return origData, err
}
