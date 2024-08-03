package message_store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb  *redis.Client
	host = "redis:6379"
)

type User struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
	validateRedis()
}

func validateRedis() {
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("redis: connected", pong)
}

// TODO Implement create session
func CreateUserSession() error {
	return nil
}

// TODO Implement auth

// TODO Publish message

func CreateUser(user User) error {
	ctx := context.Background()
	userKey := "user:" + user.Username

	//TODO check if user exists

	_, err := rdb.HMSet(ctx, userKey, map[string]interface{}{
		"username":      user.Username,
		"password_hash": user.PasswordHash,
	}).Result()

	if err != nil {
		return fmt.Errorf("Failed to save user %w", err)
	}
	return nil

}
