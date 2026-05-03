package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = time.Second * 10

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{
		rdb,
	}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (s *RedisStore) Book(b Booking) error {
	b, err := s.hold(b)

	if err != nil {
		return err
	}

	log.Printf("Session booked %v", b)

	return nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	return []Booking{}
}

func (s *RedisStore) hold(b Booking) (Booking, error) {
	b.ID = uuid.New().String()
	b.ExpiresAt = time.Now().Add(defaultHoldTTL)
	b.Status = "held"

	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
	val, _ := json.Marshal(b)

	if res := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: string(redis.NX),
		TTL:  defaultHoldTTL,
	}); res.Val() != "OK" {
		return Booking{}, ErrSeatAlreadyBooked
	}

	s.rdb.Set(ctx, sessionKey(b.ID), key, defaultHoldTTL)

	return b, nil
}

// func (s *RedisStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {

// }
// func (s *RedisStore) Release(ctx context.Context, sessionID string, userID string) error {

// }
