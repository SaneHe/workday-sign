package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type AesCbc struct {
	key []byte
	iv  []byte
}

func NewAesCbc(key, iv string) *AesCbc {
	return &AesCbc{
		key: []byte(key),
		iv:  []byte(iv),
	}
}

func (ac *AesCbc) Encrypt(data string) (string, error) {
	block, err := aes.NewCipher(ac.key)
	if err != nil {
		return "", err
	}

	content := ac.pkcs5Padding([]byte(data), block.BlockSize())
	blockMode, encrypted := cipher.NewCBCEncrypter(block, ac.iv), make([]byte, len(content))
	blockMode.CryptBlocks(encrypted, content)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// 填充 补全码
func (ac *AesCbc) pkcs5Padding(data []byte, blockSize int) []byte {
	paddingLen := blockSize - len(data)%blockSize
	paddingText := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)

	return append(data, paddingText...)
}

func (ac *AesCbc) Decrypt(data string) (string, error) {
	block, err := aes.NewCipher(ac.key)
	if err != nil {
		return "", err
	}

	content, err := base64.StdEncoding.DecodeString(data)
	if len(content) == 0 || err != nil {
		return "", errors.New("invalid data")
	}

	blockMode, decrypted := cipher.NewCBCDecrypter(block, ac.iv), make([]byte, len(content))
	blockMode.CryptBlocks(decrypted, content)

	return string(ac.pkcs5UnPadding(decrypted)), nil
}

// 去除 补全码
func (ac *AesCbc) pkcs5UnPadding(content []byte) []byte {
	length := len(content)
	return content[:(length - int(content[length-1]))]
}
