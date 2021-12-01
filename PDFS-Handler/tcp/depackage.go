package tcp

func depackage(byteStream []byte, pg *Package) error {
	pg.Op = string(byteStream[0])
	usernameLength := int(byteStream[1])

	pg.username = string(byteStream[2 : 2+usernameLength])
	if pg.Op == NEW_USER_OP || pg.Op == DEL_USER_OP ||pg.Op == LOGIN_OP {
		passwdLength :=  int(byteStream[2+usernameLength])
		pg.passwd = string(byteStream[3+usernameLength : 3+usernameLength+passwdLength])
	} else if pg.Op == CHANGE_PASSWD_OP {
		passwdLength :=  int(byteStream[2+usernameLength])
		pg.passwd = string(byteStream[3+usernameLength : 3+usernameLength+passwdLength])
		newpasswdLength := int(byteStream[3+usernameLength+passwdLength])
		pg.newpasswd = string(byteStream[4+usernameLength+passwdLength:4+usernameLength+passwdLength+newpasswdLength])
	} else if pg.Op == WRITE_OP || pg.Op == READ_OP || pg.Op == DEL_OP {
		filenameLength := int(byteStream[2+usernameLength])
		pg.filename = string(byteStream[3+usernameLength : 3+usernameLength+filenameLength])
		pathLength:= int(byteStream[3+usernameLength+filenameLength])
		pg.newpasswd = string(byteStream[4+usernameLength+filenameLength:4+usernameLength+filenameLength+pathLength])
	} else if pg.Op == NEW_PATH_OP || pg.Op == DEL_PATH_OP || pg.Op == ASK_FILES_OP {
		pathLength:= int(byteStream[2+usernameLength])
		pg.path = string(byteStream[3+usernameLength : 3+usernameLength+pathLength])
	}
	return nil
}
