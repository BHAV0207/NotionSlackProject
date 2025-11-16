package websockets

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisBroker(addr, password string) *RedisBroker {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("❌ Redis connection failed:", err)
	}

	fmt.Println("✅ Connected to Redis!")
	return &RedisBroker{Client: rdb, Ctx: ctx}
}

func (rb *RedisBroker) Publish(room string, message []byte) {
	channel := fmt.Sprintf("room:%s", room)
	rb.Client.Publish(rb.Ctx, channel, message)
}

func (rb *RedisBroker) Subscribe(room string, hub *Hub) {
	channel := fmt.Sprintf("room:%s", room)
	sub := rb.Client.Subscribe(rb.Ctx, channel)
	ch := sub.Channel()

	go func() {
		defer sub.Close()
		for msg := range ch {
			if msg == nil {
				continue
			}
			roomObj := hub.GetRoom(room)
			roomObj.Broadcast <- []byte(msg.Payload)
		}
	}()
}
