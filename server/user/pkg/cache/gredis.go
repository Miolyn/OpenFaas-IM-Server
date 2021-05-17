package cache

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

func Set(key string, data interface{}, time int) error {
	conn := Cli

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

func Exists(key string) bool {
	conn := Cli

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := Cli

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func LPushList(key string, data interface{}) error {
	conn := Cli

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("LPUSH", key, value)
	if err != nil {
		return err
	}

	return nil
}

func RPushList(key string, data interface{}) error {
	conn := Cli

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("RPUSH", key, value)
	if err != nil {
		return err
	}

	return nil
}

func GetListLen(key string) (int64, error) {
	conn := Cli

	reply, err := conn.Do("LLEN", key)
	if err != nil {
		return 0, err
	}

	return reply.(int64), nil
}

func GetLIndexList(key string, index int64) ([]byte, error) {
	conn := Cli

	reply, err := redis.Bytes(conn.Do("LINDEX", key, index))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := Cli

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := Cli

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
