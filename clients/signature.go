package clients

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateSig(message, secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		Error := Error{Code: 401, Reason: "Bad Secret Key"}
		return "", error(Error)
	}

	signature := hmac.New(sha256.New, key)
	_, err = signature.Write([]byte(message))
	if err != nil {
		Error := Error{Code: 401, Reason: "Bad API Key/Pass Phrase"}
		return "", error(Error)
	}

	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}
