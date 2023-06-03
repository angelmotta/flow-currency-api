package exchangestore

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
)

type ExchangeStore struct {
	rdb *redis.Client
}

var ctx = context.Background()

func New(addr string) (*ExchangeStore, error) {
	if addr == "" {
		log.Println("addr can not be blank")
		return nil, errors.New("DB Addr can not be blank")
	}
	rdbClient := redis.NewClient(&redis.Options{
		Addr:     addr, // Addr: "localhost:6379"
		Password: "",
		DB:       0,
	})

	// Test connection
	err := rdbClient.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &ExchangeStore{
		rdb: rdbClient,
	}, nil
}

// GetExchange retrieves a Currency Exchange value from the DB layer
func (e *ExchangeStore) GetExchange(key string) (string, error) {
	val, err := e.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("key does not exist")
		return "", nil
	} else if err != nil {
		log.Println("Redis got error in Get Operation: ", err)
		return "", err
	}
	return val, err
}

// SetExchange set a Currency Exchange value to the DB layer
func (e *ExchangeStore) SetExchange(key string, val string) error {
	err := e.rdb.Set(ctx, key, val, 0).Err()
	return err
}
