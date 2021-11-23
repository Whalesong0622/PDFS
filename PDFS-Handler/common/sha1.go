package common

import (
	"crypto/sha1"
	"encoding/hex"
)

func ToSha (path string) string {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(path))
	Result := Sha1Inst.Sum([]byte(""))
	shaString := hex.EncodeToString(Result)
	return shaString
}
