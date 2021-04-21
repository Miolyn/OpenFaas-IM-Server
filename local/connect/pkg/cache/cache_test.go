package cache

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	// cl1 := NewRedisCli(REDIS_SERVER)
	res, err := GetUserConn("test")
	fmt.Printf("%#v\n", res)
	t.Log(res, err)

}
