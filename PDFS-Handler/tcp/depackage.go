package tcp

import (
	"PDFS-Handler/cookies"
	"PDFS-Handler/errorcode"
	"fmt"
	"log"
)

// 解包只做基本的参数合法校验，如参数是否为空，长度是否合法,cookie是否存在
// 检查参数正确性应当在api中处理
func depackage(byteStream []byte, pg *Package) byte {
	pg.Op = byteStream[0]
	if pg.Op == LOGIN_OP || pg.Op == NEW_USER_OP {
		// LOGIN_IN和NEW_USER包长度应该为41，其中0为op，1-20为用户名，21-40为密码
		if len(byteStream) != 41 {
			log.Println("Error when depackage,length of package not equal to 41.")
			return errorcode.UNKNOWN_ERR
		}
		pg.username = GetPackageString(byteStream[1:21])
		pg.passwd = GetPackageString(byteStream[21:41])
		if len(pg.username) == 0 {
			log.Println("Error when depackage,username nil.")
			return errorcode.UNKNOWN_ERR
		} else if len(pg.passwd) == 0 {
			log.Println("Error when depackage,password nil.")
			return errorcode.UNKNOWN_ERR
		}

		if pg.Op == LOGIN_OP {
			log.Println("Receive login request.Username:", pg.username)
		} else if pg.Op == NEW_USER_OP {
			log.Println("Receive new user request.Username:", pg.username)
		}
		return errorcode.OK
	} else if pg.Op == DEL_USER_OP {
		// DEL_USER包长度应该为41，其中0为op，1-20为COOKIE，21-40为密码
		// 删除用户操作只能在已登录状态下才能进行，且需要再次输入密码
		if len(byteStream) != 41 {
			log.Println("Error when depackage,length of package not equal to 41.")
			return errorcode.UNKNOWN_ERR
		}
		pg.Cookie = GetPackagebytes(byteStream[1:21])
		pg.passwd = GetPackageString(byteStream[21:41])
		pg.username = cookies.CookieToUsername(pg.Cookie)
		if pg.passwd == "" {
			log.Println("Error when depackage,password nil.")
			return errorcode.UNKNOWN_ERR
		} else if pg.username == "" {
			log.Println("Error when depackage,cookie not found.")
			return errorcode.COOKIES_NOT_FOUND
		}

		log.Println("Receive delete user request.Username:", pg.username)
		return errorcode.OK
	} else if pg.Op == CHANGE_PASSWD_OP {
		// CHANGE_PASSWD包长度应该为41，其中0为op，1-20为COOKIE，21-40为新密码
		// 修改密码操作只能在已登录状态下才能进行，前端重复认证密码即可
		if len(byteStream) != 41 {
			log.Println("Error when depackage,length of package not equal to 41.")
			return errorcode.UNKNOWN_ERR
		}
		pg.Cookie = GetPackagebytes(byteStream[1:21])
		pg.passwd = GetPackageString(byteStream[21:41])
		pg.username = cookies.CookieToUsername(pg.Cookie)
		if pg.passwd == "" {
			log.Println("Error when depackage,password or cookie nil.")
			return errorcode.UNKNOWN_ERR
		} else if pg.username == "" {
			log.Println("Error when depackage,cookie not found.")
			return errorcode.COOKIES_NOT_FOUND
		}

		log.Println("Receive change password request.Username:", pg.username)
		return errorcode.OK
	} else if pg.Op == READ_OP || pg.Op == WRITE_OP || pg.Op == DEL_OP {
		if len(byteStream) <= 21 {
			// 刚好等于21时，path为空
			log.Println("Error when depackage,length of package less then 21.")
			return errorcode.UNKNOWN_ERR
		}
		pg.Cookie = GetPackagebytes(byteStream[1:21])
		pg.username = cookies.CookieToUsername(pg.Cookie)
		if pg.username == "" {
			log.Println("Error when depackage,cookie not found.")
			return errorcode.COOKIES_NOT_FOUND
		}
		pg.path = GetPackageString(byteStream[21:])

		if pg.Op == READ_OP {
			log.Println("Receive read request.Username:", pg.username, "Path:", pg.path)
		} else if pg.Op == WRITE_OP {
			log.Println("Receive write request.Username:", pg.username, "Path:", pg.path)
		} else if pg.Op == DEL_OP {
			log.Println("Receive delete request.Username:", pg.username, "Path:", pg.path)
		}
		return errorcode.OK
	} else if pg.Op == ASK_FILES_OP {
		if len(byteStream) <= 21 {
			// 刚好等于21时，path为空
			fmt.Println("Error when depackage,length of package less then 21.")
			return errorcode.UNKNOWN_ERR
		}
		pg.Cookie = GetPackagebytes(byteStream[1:21])
		pg.username = cookies.CookieToUsername(pg.Cookie)
		if pg.username == "" {
			log.Println("Error when depackage,cookie not found.")
			return errorcode.COOKIES_NOT_FOUND
		}
		pg.path = GetPackageString(byteStream[21:])

		log.Println("Receive ask files request.Username:", pg.username, "Path:", pg.path)
		return errorcode.OK
	} else if pg.Op == ADD_DIR_OP || pg.Op == DEL_DIR_OP {
		if len(byteStream) <= 21 {
			// 刚好等于21时，path为空
			fmt.Println("Error when depackage,length of package less then 21.")
			return errorcode.UNKNOWN_ERR
		}
		pg.Cookie = GetPackagebytes(byteStream[1:21])
		pg.username = cookies.CookieToUsername(pg.Cookie)
		if pg.username == "" {
			log.Println("Error when depackage,cookie not found.")
			return errorcode.COOKIES_NOT_FOUND
		}
		pg.path = GetPackageString(byteStream[22:])

		if pg.Op == ADD_DIR_OP {
			log.Println("Receive add diretory request.Username:", pg.username, "Path:", pg.path)
		} else if pg.Op == DEL_DIR_OP {
			log.Println("Receive delete diretory request.Username:", pg.username, "Path:", pg.path)
		}
		return errorcode.OK
	}
	log.Println("Operation not exist.")
	return errorcode.UNKNOWN_ERR
}

func GetPackagebytes(bytes []byte) []byte {
	// 字节流固定长度，多余部分以0x00填充，所以遇到字节0直接返回即可
	buf := make([]byte, 0)
	for _, by := range bytes {
		if by == 0 {
			break
		}
		buf = append(buf, by)
	}
	return buf
}

func GetPackageString(bytes []byte) string {
	// 字节流固定长度，多余部分以0x00填充，所以遇到字节0直接返回即可
	return string(GetPackagebytes(bytes))
}
