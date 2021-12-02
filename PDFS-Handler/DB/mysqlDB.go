package DB

import (
	"PDFS-Handler/common"
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

const (
	userName = "root"
	password = "123456"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "PDFS"
)

var createTableSQL = "CREATE TABLE if not exists `people_tb` (`username` varchar(25) DEFAULT '' UNIQUE,`passwd` varchar(50) DEFAULT '',PRIMARY KEY (`username`))ENGINE=InnoDB DEFAULT CHARSET=utf8;"

func MySQLConnect() (*sql.DB, error) {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, _ := sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		return nil, err
	}
	return DB, nil
}

func MySQLInit() {
	db, err := MySQLConnect()
	if err != nil {
		log.Println("Error occur when connecting To MySQL err:", err)
	}
	_, _ = db.Exec(createTableSQL)
}

func NewUserToDB(username string, passwd string) string {
	db, err := MySQLConnect()
	if err != nil {
		log.Println("Error occur when connecting To MySQL err:", err)
		return common.UNKNOWN_ERR
	}
	defer db.Close()

	if IsUserExist(username) {
		log.Println("Error occur when creating new user,user already exist.")
		return common.USER_EXIST
	}

	SQL := "insert into people_tb(username,passwd)values (?,?)"
	args := []string{username, common.ToSha(passwd)}

	//执行SQL语句
	_, err = db.Exec(SQL, args[0], args[1])
	if err != nil {
		log.Println("Error occur when creating new user:", err)
		return common.UNKNOWN_ERR
	} else {
		log.Println("New user", username, "created successfully.")
		return common.OK
	}
}

func DelUserToDB(username string, passwd string) string {
	db, err := MySQLConnect()
	if err != nil {
		log.Println("Error occur when connecting To MySQL err:", err)
		return common.UNKNOWN_ERR
	}
	defer db.Close()

	if !IsUserExist(username) {
		log.Println("Error occur when deleting user,", username, "not exist.")
		return common.USER_NOT_EXIST
	}

	check := PasswdCheck(username, passwd)
	if check != common.OK {
		return check
	}
	SQL := "delete from people_tb where username = ?"
	_, err = db.Exec(SQL, username)
	if err != nil {
		log.Println("Error occur when deleting user", username, ":", err)
		return common.UNKNOWN_ERR
	}
	log.Println("Delete user", username, "successfully.")
	return common.OK

}

func PasswdCheck(username string, passwd string) string {
	db, err := MySQLConnect()
	if err != nil {
		log.Println("Error occur when connecting To MySQL err:", err)
		return common.UNKNOWN_ERR
	}
	defer db.Close()

	exist := IsUserExist(username)
	if !exist {
		fmt.Println("Check user passwd failed,user", username, "not exist.")
		return common.USER_NOT_EXIST
	}
	SQL := "select * from people_tb where username = ?"
	args := []string{username}

	rows := db.QueryRow(SQL, args[0])
	var name string
	var tb_passwd string
	err = rows.Scan(&name, &tb_passwd)

	if tb_passwd == common.ToSha(passwd) {
		log.Println("Password correct:", username)
		return common.OK
	} else {
		log.Println("Password incorrect:", username)
		return common.PASSWD_ERROR
	}
}

func ChangePasswd(username string, passwd string, newpasswd string) string {
	db, err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return common.UNKNOWN_ERR
	}
	defer db.Close()

	reply := PasswdCheck(username, passwd)
	if reply != common.OK {
		return reply
	}

	SQL := "update people_tb set passwd = ? where username = ?"
	args := []string{common.ToSha(newpasswd), username}

	row, err := db.Exec(SQL, args[0], args[1])
	if err != nil {
		log.Println("Error occur when changing user passwd:", username, err)
		return common.UNKNOWN_ERR
	}

	rowsaffected, err := row.RowsAffected()
	if err != nil || rowsaffected != 1 {
		log.Println("Error occur when changing user passwd:", username, err)
		return common.UNKNOWN_ERR
	}
	log.Println("Change", username, "password successfully.")
	return common.OK
}

func IsUserExist(username string) bool {
	db, err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	rows := db.QueryRow("select * from people_tb where username = ?", username)
	if err != nil {
		return false
	}
	var name string
	var tb_passwd string
	err = rows.Scan(&name, &tb_passwd)
	if name == "" {
		return false
	} else {
		return true
	}
}
