package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/gofrs/uuid"
)

func GenerateId() string {
	uuid, _ := uuid.NewV4()
	md5Bytes := md5.Sum(uuid.Bytes())
	md5Str := fmt.Sprintf("%x", md5Bytes)
	return md5Str[0:15]
}

func GenerateSHA256(str string) []byte {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashed := hash.Sum(nil)
	return hashed
}
