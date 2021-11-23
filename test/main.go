package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func main() {
	TestString := "Hi, pandaman!"
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(TestString))
	Result := Sha1Inst.Sum([]byte(""))
	fmt.Printf("%x\n", Result)
	mystring := hex.EncodeToString(Result)
	fmt.Println(mystring)
}

