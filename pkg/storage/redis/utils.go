package redis

import (
	"crypto/md5"
	"encoding/hex"
)

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text)) //nolint
	return hex.EncodeToString(hasher.Sum(nil))
}
