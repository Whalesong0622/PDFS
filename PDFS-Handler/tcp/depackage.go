package tcp

import (
	"log"
	"strconv"
)

func depackage(byteStream []byte, pg *Package) error {
	pg.Op = string(byteStream[0])
	usernameLength, err := strconv.Atoi(string(byteStream[1]))
	if err != nil {
		log.Println("Error occur when depackaging:", err)
		return err
	}
	pg.username = string(byteStream[2 : 2+usernameLength])
	if pg.Op == NEW_USER_OP || pg.Op == DEL_USER_OP ||pg.Op == LOGIN_OP{
		passwdLength, err := strconv.Atoi(string(byteStream[2+usernameLength : 3+usernameLength]))
		if err != nil {
			log.Println("Error occur when depackaging:", err)
			return err
		}
		pg.passwd = string(byteStream[3+usernameLength : 3+usernameLength+passwdLength])
	} else if pg.Op == CHANGE_PASSWD_OP {
		passwdLength, err := strconv.Atoi(string(byteStream[2+usernameLength : 3+usernameLength]))
		if err != nil {
			log.Println("Error occur when depackaging:", err)
			return err
		}
		pg.passwd = string(byteStream[3+usernameLength : 3+usernameLength+passwdLength])
		newpasswdLength,err := strconv.Atoi(string(byteStream[3+usernameLength+passwdLength : 4+usernameLength+passwdLength]))
		if err != nil {
			log.Println("Error occur when depackaging:", err)
			return err
		}
		pg.newpasswd = string(byteStream[4+usernameLength+passwdLength:4+usernameLength+passwdLength+newpasswdLength])
	} else if pg.Op == WRITE_OP || pg.Op == READ_OP || pg.Op == DEL_OP {
		filenameLength, err := strconv.Atoi(string(byteStream[2+usernameLength : 3+usernameLength]))
		if err != nil {
			log.Println("Error occur when depackaging:", err)
			return err
		}
		pg.filename = string(byteStream[3+usernameLength : 3+usernameLength+filenameLength])
		pathLength,err := strconv.Atoi(string(byteStream[3+usernameLength+filenameLength : 4+usernameLength+filenameLength]))
		if err != nil {
			log.Println("Error occur when depackaging:", err)
			return err
		}
		pg.newpasswd = string(byteStream[4+usernameLength+filenameLength:4+usernameLength+filenameLength+pathLength])
	} else if pg.Op == NEW_PATH_OP || pg.Op == DEL_PATH_OP || pg.Op == ASK_FILES_OP {
		pathLength, err := strconv.Atoi(string(byteStream[2+usernameLength : 3+usernameLength]))
		if err != nil {
			log.Println("Error occur when depackaging:", err)
			return err
		}
		pg.path = string(byteStream[3+usernameLength : 3+usernameLength+pathLength])
	}
	return nil
}
