package other

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	//	hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func Md5ByString(str string) string {
	// 方法一
	//data := []byte(str)
	//has := md5.Sum(data)
	//fmt.Sprintf("%x", has)

	// 方法二
	m := md5.New()
	_, _ = io.WriteString(m, str)
	arr := m.Sum(nil)
	// hex.EncodeToString(arr)
	return fmt.Sprintf("%x", arr)
}
