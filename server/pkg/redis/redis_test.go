package redis

import (
	"testing"

	"github.com/go-redis/redismock/v9"
)

var rdb *RedisClient

func TestSetAndGet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	rdb = Initialize(db)
}
