package main

import (
	"crypto/sha1"
	"database/sql"
	_ "database/sql"
	"encoding/hex"
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

func ToSha(path string) string {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(path))
	Result := Sha1Inst.Sum([]byte(""))
	shaString := hex.EncodeToString(Result)
	return shaString
}

func NewUserToDB(username string, passwd string) error {
	db, err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	SQL := "insert into people_tb(username,passwd)values (?,?)"
	args := []string{username, ToSha(passwd)}

	//执行SQL语句
	r, err := db.Exec(SQL, args[0], args[1])
	if err != nil {
		log.Println("Error occur when insert new user:", err)
		return err
	}

	//查询最后一天用户ID，判断是否插入成功
	_, err = r.RowsAffected()
	if err != nil {
		log.Println("Error occur when insert new user:", err)
		return err
	}
	log.Println("New user", username, "created successfully.")
	return nil
}

func DelUserToDB(username string, passwd string) error {
	db, err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	rows := db.QueryRow("select * from people_tb where username = ?", username)
	if err != nil {
		return err
	}
	var name string
	var tb_passwd string
	err = rows.Scan(&name, &tb_passwd)
	fmt.Println(err)
	fmt.Println(name, tb_passwd)
	if tb_passwd == ToSha(passwd) {
		_, err := db.Exec("delete from people_tb where username = ?", username)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		fmt.Println("Passwd wrong")
	}
	return nil
}

func LoginPasswdCheck(username string, passwd string) bool {
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
	fmt.Println(err)
	fmt.Println(name, tb_passwd)
	if name == ""{
		log.Printf("haha")
	}
	if tb_passwd == ToSha(passwd) {
		return true
	} else {
		fmt.Println("Passwd wrong")
		return false
	}
	return false
}

var createTableSQL = "CREATE TABLE if not exists `people_tb` (`username` varchar(25) DEFAULT '' UNIQUE,`passwd` varchar(50) DEFAULT '',PRIMARY KEY (`username`))ENGINE=InnoDB DEFAULT CHARSET=utf8;"
func main() {
	log.Println("1新建2删除3登陆")
	db,err := MySQLConnect()
	if err != nil {
		fmt.Println(err)
		return
	}
	reply,err := db.Exec(createTableSQL)
	fmt.Println(reply,err)
	for {
		var op string
		_, _ = fmt.Scan(&op)
		var username string
		var passwd string
		_, _ = fmt.Scan(&username)
		_, _ = fmt.Scan(&passwd)
		fmt.Println(ToSha(passwd))
		fmt.Println(op, username, passwd)
		if op == "1" {
			_ = NewUserToDB(username, passwd)
		} else if op == "2" {
			_ = DelUserToDB(username, passwd)
		} else if op == "3" {
			_ = LoginPasswdCheck(username, passwd)
		}
	}
}

/*
CREATE TABLE `people_tb` (
    `username` varchar(25) DEFAULT '' UNIQUE,
    `passwd` varchar(50) DEFAULT '',
    PRIMARY KEY (`username`)
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/
