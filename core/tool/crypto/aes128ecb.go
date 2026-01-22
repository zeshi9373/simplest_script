package crypto

import (
	"bytes"
	"crypto/aes"
	"fmt"
)

func pad(data []byte) []byte {
	padLen := aes.BlockSize - len(data)%aes.BlockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func unpad(data []byte) ([]byte, error) {
	length := len(data)
	unpadLen := int(data[length-1])
	if unpadLen > length {
		return nil, fmt.Errorf("unpad error")
	}
	return data[:(length - unpadLen)], nil
}

func AesECBEncrypt(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Padding the plaintext
	paddedText := pad(plainText)

	encrypted := make([]byte, len(paddedText))
	for i := 0; i < len(paddedText); i += aes.BlockSize {
		block.Encrypt(encrypted[i:i+aes.BlockSize], paddedText[i:i+aes.BlockSize])
	}
	return encrypted, nil
}
