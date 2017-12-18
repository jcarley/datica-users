package encoding

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"hash/adler32"
	"io"
	"log"

	"golang.org/x/crypto/scrypt"
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

type Hasher struct {
	Algorithm string
	Buffer    *bytes.Buffer
	Hash      hash.Hash
}

func GenerateSalt() string {
	salt := make([]byte, PW_SALT_BYTES)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", salt)
}

func Key(password, salt string) string {
	hash, err := scrypt.Key([]byte(password), []byte(salt), 1<<14, 8, 1, PW_HASH_BYTES)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", hash)
}

func HashPassword(password string) string {
	b := bytes.NewBufferString(password)
	h := NewHasher("SHA256", b)
	return h.HashString()
}

func NewHasher(algorithm string, buffer *bytes.Buffer) *Hasher {
	hasher := Hasher{Algorithm: algorithm, Buffer: buffer}
	hasher.initAlgorithm()
	return &hasher
}

func (h *Hasher) initAlgorithm() {
	switch h.Algorithm {
	case "SHA1":
		h.Hash = sha1.New()
	case "SHA256":
		h.Hash = sha256.New()
	case "MD5":
		h.Hash = md5.New()
	case "ADLER32":
		h.Hash = adler32.New()
	}
}

func (h *Hasher) HashString() string {
	h.Hash.Write(h.Buffer.Bytes())
	sum := h.Hash.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)
}
