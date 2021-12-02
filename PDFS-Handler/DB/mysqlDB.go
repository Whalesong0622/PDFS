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
	port     = "3380"
	dbName   = "mysql"
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
	db,err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, _ = db.Exec(createTableSQL)
}


func NewUserToDB(username string, passwd string) bool {
	db, err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	SQL := "insert into people_tb(username,passwd)values (?,?)"
	args := []string{username, common.ToSha(passwd)}

	//执行SQL语句
	_, err = db.Exec(SQL, args[0], args[1])
	if err != nil {
		log.Println("Error occur when creating new user:", err)
		return false
	}else{
		log.Println("New user", username, "created successfully.")
		return true
	}
}

func DelUserToDB(username string, passwd string) bool {
	db, err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()

	SQL := "select * from people_tb where username = ?"
	args := []string{username, common.ToSha(passwd)}

	rows := db.QueryRow(SQL, args[0], args[1])
	if err != nil {
		log.Println("Error occur when deleting user:", username,err)
		return false
	}
	var name string
	var tb_passwd string
	err = rows.Scan(&name, &tb_passwd)

	if tb_passwd == common.ToSha(passwd) {
		SQL = "delete from people_tb where username = ?"
		_, err := db.Exec(SQL, username)
		if err != nil {
			log.Println("Error occur when deleting user:",username, err)
			return false
		}
		log.Println("Delect user",username,"successfully.")
		return true
	} else {
		log.Println("Delect user",username,"failed,password incorrect.")
		return false
	}
}

func LoginPasswdCheck(username string, passwd string) string {
	db, err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return common.UNKNOWN_ERR
	}
	defer db.Close()

	exist := IsUserExist(username)
	if !exist {
		fmt.Println("Create user failed,user not exist.")
		return common.USER_EXIST
	}
	SQL := "select * from people_tb where username = ?"
	args := []string{username, common.ToSha(passwd)}

	rows := db.QueryRow(SQL, args[0], args[1])
	var name string
	var tb_passwd string
	err = rows.Scan(&name, &tb_passwd)

	if tb_passwd == common.ToSha(passwd) {
		log.Println("Password correct,",username,"login successfully.")
		return common.OK
	} else {
		log.Println("Password incorrect,",username,"login failed.")
		return common.PASSWD_ERROR
	}
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
	}else{
		return true
	}
}