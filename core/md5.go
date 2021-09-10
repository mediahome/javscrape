package core

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(value []byte) string {
	md := md5.Sum(value)
	return hex.EncodeToString(md[:])
}
