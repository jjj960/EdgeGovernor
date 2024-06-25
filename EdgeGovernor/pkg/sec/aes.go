package sec

import (
	"crypto/aes"
	"crypto/cipher"
	"log"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
)

const (
	blockSize        = 256 / 8
	pbkdf2Iterations = 4096
	salt             = "EDGEGOVER-NOR-RANDOM-AND-SECRET-SALT"
)

var iv [blockSize / 2]byte

type AES struct {
	block cipher.Block
}

var Safer *AES

func NewMSecLayer(pass string) *AES {
	key := pbkdf2.Key([]byte(pass), []byte(salt), pbkdf2Iterations, blockSize, sha3.New256)
	if len(key) != blockSize {
		log.Panicf("key.len = %d and blocksize = %d, mismatch", len(key), blockSize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic("failed to create aes cipher: ", err)
	}

	return &AES{block: block}
}

func (sec *AES) Encrypt(in []byte) []byte { //CFB模式加密
	out := make([]byte, len(in))
	cipher.NewCFBEncrypter(sec.block, iv[:]).XORKeyStream(out, in)
	return out
}

func (sec *AES) Decrypt(in []byte) []byte {
	out := make([]byte, len(in))
	cipher.NewCFBDecrypter(sec.block, iv[:]).XORKeyStream(out, in)
	return out
}
