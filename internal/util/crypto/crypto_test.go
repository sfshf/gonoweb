package crypto_test

import (
	"testing"

	"github.com/sfshf/gonoweb/internal/util/crypto"
	"github.com/stretchr/testify/assert"
)

// go test -v -count=1 -timeout 30s -run ^TestMd5Hex$ github.com/sfshf/gonoweb/internal/util/crypto
func TestMd5Hex(t *testing.T) {
	encrypted := crypto.Md5Hex("123")
	assert.Equal(t, encrypted, "202cb962ac59075b964b07152d234b70")
}
