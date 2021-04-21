package cache

import (
	"OpenFaas-Connect/pkg/logger"
	"OpenFaas-Connect/st"
	"fmt"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

var (
	REDIS_SERVER = "49.4.114.179:6378"
	cli          redis.Conn
)

func init() {
	cli = NewRedisCli(REDIS_SERVER)

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

func SetUser(userId string, connID int64) error {
	key := getUserKey(userId)
	_, err := cli.Do("SET", key, connID)
	return err
}

func DelUser(userId string) error {
	key := getUserKey(userId)
	_, err := cli.Do("DEL", key)
	return err

}

func GetUserConn(userId string) (int64, error) {
	key := getUserKey(userId)
	res, err := cli.Do("GET", key)
	if res == nil || err != nil {
		return 0, fmt.Errorf("not found key :%v  err:%v", key, err)
	}
	connId, _ := strconv.ParseInt(fmt.Sprintf("%s", res), 10, 64)
	return connId, nil
}

func getUserKey(userId string) string {
	return fmt.Sprintf("uid_%v", userId)
}
