package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
)

func GenerateRequestId() string {
	uuid := uuid.NewString()
	md5Bytes := md5.Sum([]byte(uuid))
	md5Str := fmt.Sprintf("%x", md5Bytes)
	return md5Str[0:15]
}

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func GenerateKeyId() string {
	return uuid.NewString()
}

func GenerateSHA256(str string) []byte {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashed := hash.Sum(nil)
	return hashed
}
