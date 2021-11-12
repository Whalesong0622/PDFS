package mysqldb

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	userName = "root"
	password = "123456"
	ip       = "1.117.64.144"
	port     = "3380"
	dbName   = "mysql"
)

func InitDB() *sql.DB {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, _ := sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		return nil
	}
	fmt.Println("connnect success")
	return DB
}
