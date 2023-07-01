package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"regexp"
)

var urlRegex = regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[.\!\/\\w]*))?)`)

func MakeMD5UUID() string {
	hash := md5.Sum([]byte(uuid.NewString()))
	return hex.EncodeToString(hash[:])
}

func ParseURL(url string) bool {
	return urlRegex.MatchString(url)
}
