package tool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

func Md5(txt string) string {
	m5 := md5.New()
	m5.Write([]byte(txt))
	txtHash := hex.EncodeToString(m5.Sum(nil))
	return txtHash
}

// HmacSha256 计算HmacSha256
// key 是加密所使用的key
// data 是加密的内容
func HmacSha256(key string, data string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(data))

	return mac.Sum(nil)
}

// HmacSha256ToHex 将加密后的二进制转16进制字符串
func HmacSha256ToHex(key string, data string) string {
	return strings.ToLower(hex.EncodeToString(HmacSha256(key, data)))
}

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// AES CBC 模式加密
func AesEncrypt(plainText, key []byte, isTransfer bool) (string, error) {
	// 生成 AES 块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7 填充
	plainText = pkcs7Padding(plainText, block.BlockSize())

	// 生成随机 IV
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// CBC 加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	enstr := base64.StdEncoding.EncodeToString(cipherText)

	// 返回 base64 编码后的密文
	if isTransfer {
		return strings.ReplaceAll(enstr, "+", "_"), nil
	} else {
		return enstr, nil
	}
}

// PKCS7 去除填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("data length is zero")
	}

	padding := int(data[length-1])
	if padding > length || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding size")
	}

	return data[:length-padding], nil
}

// AES CBC 模式解密
func AesDecrypt(cipherTextBase64 string, key []byte, isTransfer bool) (string, error) {
	if isTransfer {
		cipherTextBase64 = strings.ReplaceAll(cipherTextBase64, "_", "+")
	}

	// 解码 base64 密文
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}

	// 生成 AES 块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("cipherText too short")
	}

	// 读取 IV
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	// CBC 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	// 去除 PKCS7 填充
	plainText, err := pkcs7Unpadding(cipherText)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func Sha256ToHex(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func AesEcbEncrypt(plainText, base64Key string) string {
	// 解码base64 key
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return ""
	}

	// 创建AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	// ECB模式需要手动处理块对齐
	data := []byte(plainText)
	blockSize := block.BlockSize()

	// PKCS7填充
	padText := pkcs7Padding(data, block.BlockSize())
	// ECB模式加密
	encrypted := make([]byte, len(padText))
	for i := 0; i < len(padText); i += blockSize {
		block.Encrypt(encrypted[i:i+blockSize], padText[i:i+blockSize])
	}

	return base64.StdEncoding.EncodeToString(encrypted)
}
