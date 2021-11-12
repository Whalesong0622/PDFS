package common

import (
	"crypto/sha256"
	"encoding/hex"
)

func ToSha(path string) string {
	ShaInst := sha256.New()
	ShaInst.Write([]byte(path))
	Result := ShaInst.Sum([]byte(""))
	shaString := hex.EncodeToString(Result)
	return shaString
}
