package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

const (
	AES_KEY = "6e16ae7c17b32ba0511262a34202a594"
)

// 极简AES CBC加密
func AESCbcEncrypt(plainText, key string) (string, error) {
	// 补齐密钥为32字节
	paddedKey := padKey(key)

	block, err := aes.NewCipher(paddedKey)
	if err != nil {
		return "", err
	}

	// 填充明文
	paddedText := pkcs7Pad([]byte(plainText), aes.BlockSize)

	// 生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 加密
	cipherText := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, paddedText)

	// 返回 iv + cipherText 的base64
	result := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(result), nil
}

// 极简AES CBC解密
func AESCbcDecrypt(encrypted, key string) (string, error) {
	// 补齐密钥为32字节
	paddedKey := padKey(key)

	block, err := aes.NewCipher(paddedKey)
	if err != nil {
		return "", err
	}

	// 解码base64
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", err
	}

	// 提取IV和密文
	iv := data[:aes.BlockSize]
	cipherText := data[aes.BlockSize:]

	// 解密
	plainText := make([]byte, len(cipherText))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, cipherText)

	// 去除填充
	result, err := pkcs7Unpad(plainText)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// 简单的密钥补齐
func padKey(key string) []byte {
	keyBytes := []byte(key)

	// 如果密钥长度不对，简单补齐
	if len(keyBytes) >= 32 {
		return keyBytes[:32]
	}

	// 简单复制补齐
	padded := make([]byte, 32)
	copy(padded, keyBytes)
	for i := len(keyBytes); i < 32; i++ {
		padded[i] = byte(i)
	}
	return padded
}

// PKCS7填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// PKCS7去除填充
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data length is zero")
	}

	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New("invalid padding size")
	}

	return data[:len(data)-padding], nil
}
