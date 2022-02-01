package DB

import (
	"PDFS-Handler/common"
	"PDFS-Handler/errorcode"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var MySQLConfig *common.MySQLConfigStruct

var createDatabaseSQL = "CREATE DATABASE IF NOT EXISTS ?;"
var createTableSQL = "CREATE TABLE if not exists ? (`username` varchar(25) DEFAULT '' UNIQUE,`passwd` varchar(80) DEFAULT '',PRIMARY KEY (`username`))ENGINE=InnoDB DEFAULT CHARSET=utf8;"

func MySQLInit() {
	MySQLConfig = common.GetMySQLStruct()
	db, err := MySQLConnect()
	if err != nil {
		log.Println("Error occur when connecting to MySQL:", err)
	}
	_, _ = db.Exec(createDatabaseSQL, MySQLConfig.DBName)
	_, _ = db.Exec(createTableSQL, MySQLConfig.TableName)
}

func MySQLConnect() (*sql.DB, error) {
	path := strings.Join([]string{MySQLConfig.Username, ":", MySQLConfig.Passwd, "@tcp(", MySQLConfig.Ip, ":", MySQLConfig.Port, ")/", MySQLConfig.DBName, "?charset=utf8"}, "")
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

func NewUserToDB(username string, passwd string) byte {
	db, err := MySQLConnect()
	if err != nil {
		log.Println("Error occur when connecting to MySQL err:", err)
		return errorcode.UNKNOWN_ERR
	}
	defer db.Close()

	if IsUserExist(username) {
		log.Println("Error occur when creating new user,user already exist.")
		return errorcode.USER_EXIST
	}

	//执行SQL语句
	SQL := "insert into PDFS_USER_TABLE(username,passwd)values (?,?)"
	_, err = db.Exec(SQL, username, common.ToSha(passwd))
	if err != nil {
		log.Println("Error occur when executive new user:", err)
		return errorcode.UNKNOWN_ERR
	} else {
		log.Println("New user", username, "create to MySQL success.")
		return errorcode.OK
	}
}

func DelUserToDB(username string, passwd string) byte {
	db, err := MySQLConnect()
	if err != nil {
		log.Println("Error occur when connecting To MySQL err:", err)
		return errorcode.UNKNOWN_ERR
	}
	defer db.Close()

	if !IsUserExist(username) {
		log.Println("Error occur when deleting user,", username, "not exist.")
		return errorcode.USER_NOT_EXIST
	}

	check := PasswdCheck(username, passwd)
	if check != errorcode.OK {
		return check
	}
	SQL := "delete from PDFS_USER_TABLE where username = ?"
	_, err = db.Exec(SQL, username)
	if err != nil {
		log.Println("Error occur when deleting user", username, ":", err)
		return errorcode.UNKNOWN_ERR
	}
	log.Println("Delete user", username, "successfully.")
	return errorcode.OK

}

func PasswdCheck(username string, passwd string) byte {
	db, err := MySQLConnect()
	if err != nil {
		log.Println("Error occur when connecting To MySQL err:", err)
		return errorcode.UNKNOWN_ERR
	}
	defer db.Close()

	exist := IsUserExist(username)
	if !exist {
		fmt.Println("Check user passwd failed,user", username, "not exist.")
		return errorcode.USER_NOT_EXIST
	}
	SQL := "select * from PDFS_USER_TABLE where username = ?"
	args := []string{username}

	rows := db.QueryRow(SQL, args[0])
	var name string
	var tb_passwd string
	err = rows.Scan(&name, &tb_passwd)

	if err != nil {
		return errorcode.UNKNOWN_ERR
	}
	if tb_passwd == common.ToSha(passwd) {
		log.Println("Password correct:", username)
		return errorcode.OK
	} else {
		log.Println("Password incorrect:", username)
		return errorcode.PASSWD_ERROR
	}
}

func ChangePasswd(username string, newpasswd string) byte {
	db, err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return errorcode.CHANGE_PASSWD_FAILED
	}
	defer db.Close()

	SQL := "update PDFS_USER_TABLE set passwd = ? where username = ?"
	args := []string{common.ToSha(newpasswd), username}

	row, err := db.Exec(SQL, args[0], args[1])
	if err != nil {
		log.Println("Error occur when changing user passwd:", username, err)
		return errorcode.CHANGE_PASSWD_FAILED
	}

	rowsaffected, err := row.RowsAffected()
	if err != nil || rowsaffected != 1 {
		log.Println("Error occur when changing user passwd,password is same as before.", username, err)
		return errorcode.CHANGE_PASSWD_FAILED
	}
	log.Println("Change", username, "password successfully.")
	return errorcode.CHANGE_PASSWD_SUCCESS
}

func IsUserExist(username string) bool {
	db, err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	rows := db.QueryRow("select * from PDFS_USER_TABLE where username = ?", username)
	if err != nil {
		return false
	}
	var name string
	var tb_passwd string
	err = rows.Scan(&name, &tb_passwd)
	if name == "" || err != nil {
		return false
	} else {
		return true
	}
}
