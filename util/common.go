package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
)

func MakeMD5UUID() string {
	hash := md5.Sum([]byte(uuid.NewString()))
	return hex.EncodeToString(hash[:])
}
