package secret_test

import (
	"encoding/base64"
	"testing"

	"github.com/infraboard/mcube/v2/crypto/cbc"
)

// 基于tdd的方式产生SECURE_KEY
func TestNewSecret(t *testing.T) {
	t.Log(base64.StdEncoding.EncodeToString(cbc.MustGenRandomKey(cbc.AES_KEY_LEN_32)))
}
