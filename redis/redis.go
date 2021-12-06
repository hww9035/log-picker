package redis

import (
	"github.com/gomodule/redigo/redis"
)

var client redis.Conn

func Init(address string) (err error) {
	client, err = redis.Dial("tcp", address)
	if err != nil {
		return err
	}
	return err
}

func GetClient() redis.Conn {
	return client
}

func CloseRedisConn(conn redis.Conn) {
	conn.Close()
}

func Set(key string, value string) bool {
	defer CloseRedisConn(client)
	_, err := client.Do("SET", key, value)
	return err == nil
}

func Get(key string) (string, error) {
	defer CloseRedisConn(client)
	s, err := redis.String(client.Do("GET", key))
	return s, err
}
