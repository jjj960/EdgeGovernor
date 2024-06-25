package sec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
)

func TestNewMSecLayer(t *testing.T) {
	sec := NewMSecLayer("ASDASDASZXasdzxczxgg")

	strs := []string{
		"你好",
		"abcde",
		"%%%%",
		"ad56",
	}

	for i := 0; i < len(strs); i++ {
		encrpy := sec.Encrypt([]byte(strs[i]))
		s := sec.Decrypt(encrpy)
		fmt.Println("jiami:", encrpy)
		fmt.Println("jiemi:", string(s))
		if !bytes.Equal([]byte(strs[i]), s) {
			t.Errorf("coudln't encrypt #%d, input %v output %v", i, []byte(strs[i]), s)
		}

	}
}

var Key = []byte("PBWUtoFHtXeEGJyA")

func EncryptAES(plaintext string) string {
	block, err := aes.NewCipher(Key)
	if err != nil {
		return ""
	}
	in := []byte(plaintext)
	leng := len(plaintext)
	if leng%16 != 0 {
		leng = leng/16*16 + 16
		leng = leng - len(plaintext)
		for i := 0; i < leng; i++ {
			in = append(in, 0)
		}
		leng = len(in)
	}

	cipherText := make([]byte, aes.BlockSize+leng)
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(cipherText[aes.BlockSize:], in)
	return hex.EncodeToString(cipherText)
}

func DecryptAES(ct string) string {
	ciphertext, err := hex.DecodeString(ct)
	if err != nil {
		return ""
	}
	block, err := aes.NewCipher(Key)
	if err != nil {
		return ""
	}
	if len(ciphertext) < aes.BlockSize {
		return ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	// fmt.Println(len(ciphertext), len(iv))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(ciphertext, ciphertext)
	return string(ciphertext)
}
