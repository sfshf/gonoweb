package crypto

import (
	"crypto/md5"
	"fmt"
)

func Md5Hex(plain string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(plain)))
}
