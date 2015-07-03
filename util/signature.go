package util

import (
	"io"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

type Signature struct{
 	Secret string
}

func (sig Signature) Sign(message string) string {
	if len(sig.Secret) == 0 {
		return ""
	}
	mac := hmac.New(sha1.New, []byte(sig.Secret))
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}

func (sig Signature) Verify(message, signature string) bool {
	return hmac.Equal([]byte(signature), []byte(sig.Sign(message)))
}

func SignString(secret, message string) string {
	sig := Signature{secret}
	return sig.Sign(message)
}