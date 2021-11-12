package common

import "os"

func ByteToBytes(bt byte) (bytes []byte) {
	bytes = append(bytes, bt)
	return bytes
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
