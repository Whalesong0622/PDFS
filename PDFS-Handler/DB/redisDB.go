package DB

import (
	"PDFS-Handler/common"
	"github.com/gomodule/redigo/redis"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Pool *redis.Pool
)

func RedisInit() {
	master := common.GetMasterIpConfig()
	redisHost := master + ":6379"
	Pool = newPool(redisHost)
	close()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func close() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGKILL)
	go func() {
		<-signalChan
		Pool.Close()
		os.Exit(0)
	}()
}

func GetFileBlockNums(path string) (int, error) {
	conn := Pool.Get()
	defer conn.Close()

	reply, err := conn.Do("HGET", path, "blocknums")
	if err != nil {
		return 0, err
	}
	return reply.(int), err
}

func GetBlockIpList(BlockName string)([]string,error){
	conn := Pool.Get()
	defer conn.Close()

	reply, err := conn.Do("HKEYS",BlockName)
	if err != nil{
		return nil,err
	}
	return reply.([]string),err
}