package main

import (
	"PDFS-Server/api"
	"fmt"
	"os"
	_ "strconv"
	"strings"
)

func main(){
	Args := os.Args
	fmt.Println(Args)
	words := strings.FieldsFunc(Args[1],func (r rune)bool{
		return r == '/'
	})
	fmt.Println(words[len(words)-1])
	if len(Args) == 2 {
		err := api.WriteInDistributed("pdfs", Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		FileName := words[len(words)-1]
		Path := strings.Join([]string{"/Users/whaleshark/Downloads/", "pdfs", "/", FileName}, "")
		err = api.Read(Path, FileName,17)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}