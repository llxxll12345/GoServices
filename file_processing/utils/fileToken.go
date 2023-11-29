package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateToken(filename string) string {
	timestamp := time.Now().UnixNano()
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s%d", filename, timestamp)))
	fileToken := hex.EncodeToString(hash.Sum(nil))
	return fileToken
}
