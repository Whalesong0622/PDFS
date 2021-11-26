package common

import (
	"log"
)

const WRITE_OP = "1"
const READ_OP = "2"

func IsLegal(RemoteAddr string,user string,op string,path string,filename string) bool {
	if user == "" {
		log.Println("Error occur when serving",RemoteAddr,",user nil")
		return false
	}else if op == "" {
		log.Println("Error occur when serving",RemoteAddr,",operation nil")
		return false
	}else if op == WRITE_OP{
		if path == ""{
			log.Println("Error occur when serving",RemoteAddr,",Write operation but path nil")
			return false
		} else if !IsDir(path){
			log.Println("Error occur when serving",RemoteAddr,",path not exist")
			return false
		}
	}else if op == READ_OP{
		if path == ""{
			log.Println("Error occur when serving",RemoteAddr,",Write operation but path nil")
			return false
		} else if !IsFile(path+filename){
			log.Println("Error occur when serving",RemoteAddr,",file or path not exist")
			return false
		}
	}else{
		log.Println("Error occur when serving",RemoteAddr,",operation illegal")
		return false
	}
	return true
}
