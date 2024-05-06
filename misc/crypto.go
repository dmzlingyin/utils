package misc

import (
	"crypto/sha256"
	"fmt"
)

func Sha256(str string) string {
	hash := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", hash)
}
