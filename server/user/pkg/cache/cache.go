package cache

import (
	"OpenFaas-User/pkg/logger"
	"OpenFaas-User/st"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

var (
	REDIS_SERVER = "49.4.114.179:6378"
	Cli          redis.Conn
)

func init() {
	Cli = NewRedisCli(REDIS_SERVER)

}

func NewRedisCli(address string) redis.Conn {
	//cli, err := redis.Dial("tcp", REDIS_SERVER)
	cli, err := redis.Dial("tcp", address)
	if err != nil {
		logger.Logger.Fatalf("连接redis 失败 %v", err)
	}
	err = cli.Send("auth", "123456")
	if err != nil {
		fmt.Println(err)
		logger.Logger.Fatal("connect redis fail")
	}
	st.Debug("redis connect")
	return cli
}
