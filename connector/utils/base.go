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
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strconv"
	"time"
)

const (
	TimeFormat  = "2006-01-02 15:04:05"
	letterBytes = "abcdefghijkmnpqrstuvwxyz123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijkmnpqrst123456789ABCDEFGHIJKLMNPQRSTUVWXYZ"
	letLen      = int64(len(letterBytes))
)

func init() {
	rand.Seed(time.Now().Unix())
}

func IntToStr(i int64) string {
	s := strconv.FormatInt(i, 10)
	return s
}

func StrToInt(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

// s
func Timestamp() int64 {
	return time.Now().Unix()
}

// ms
func Microstamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// us
func Usstamp() int64 {
	return time.Now().UnixNano() / 1e3
}

// ns
func Nanostamp() int64 {
	return time.Now().UnixNano()
}

func HS256Sign(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(decrypted []byte) []byte {
	length := len(decrypted)
	unPadding := int(decrypted[length-1])
	return decrypted[:(length - unPadding)]
}

func AesCbcEncrypt(data, key string) string {
	encryptBytes := []byte(data)
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return ""
	}

	blockSize := block.BlockSize()
	encryptBytes = pkcs5Padding(encryptBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	encrypted := make([]byte, len(encryptBytes))
	blockMode.CryptBlocks(encrypted, encryptBytes)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func AesCbcDecript(data, key string) []byte {
	dbyte, _ := base64.StdEncoding.DecodeString(data)
	l := len(dbyte)
	if l == 0 {
		return []byte("")
	}
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return []byte("")
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	decrypted := make([]byte, l)

	blockMode.CryptBlocks(decrypted, dbyte)
	decrypted = pkcs5UnPadding(decrypted)
	return decrypted
}

func AesEcbDecrypt(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}
	return pkcs5UnPadding(decrypted)
}

func AesEcbEncrypt(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	data = pkcs5Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()
	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}
	return decrypted
}

func StrToMD5(t string) string {
	h := md5.New()
	h.Write([]byte(t))
	return hex.EncodeToString(h.Sum(nil))
}

func GetUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}

func GetSha256(data []byte) string {
	sha256Contain := sha256.New()
	sha256Contain.Write(data)
	return hex.EncodeToString(sha256Contain.Sum(nil))
}
