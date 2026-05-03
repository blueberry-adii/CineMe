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

/*
* TTL time for seat holding for a user
* Seat Released if not booked within that time
 */
const defaultHoldTTL = time.Second * 10

/*
* Store 3 (Final):
* Redis store was created to implement multiple user bookings concurrently
* Result:
* Success, as Redis is single threaded, no overrides and race conditions
 */
/*
* Depends on Redis Client
 */
type RedisStore struct {
	rdb *redis.Client
}

/*
* Redis Store Constructor
 */
func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{
		rdb,
	}
}

/*
* Returns redis key for a given session Id
 */
func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

/*
* Holds the seat for a given seat in a booking
 */
func (s *RedisStore) Book(b Booking) (Booking, error) {
	b, err := s.hold(b)

	if err != nil {
		return Booking{}, err
	}

	log.Printf("Session booked %v", b)

	return b, nil
}

/*
* Lists all the bookings from Redis for a given movieId
 */
func (s *RedisStore) ListBookings(movieID string) []Booking {
	ctx := context.Background()
	pattern := fmt.Sprintf("seat:%s:*", movieID)
	keys, err := s.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return []Booking{}
	}

	bookings := make([]Booking, 0, len(keys))
	for _, key := range keys {
		val, err := s.rdb.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var b Booking
		if err := json.Unmarshal([]byte(val), &b); err == nil {
			bookings = append(bookings, b)
		}
	}
	return bookings
}

/*
* Creates Session/Booking ID, adds expiry and sets status "held"
* saves key value pair into redis with expiry only if the pair doesnt
* already exist inside Redis.
 */
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

/*
* Confirms a booking by getting the session from session ID and user ID
* removes expiry time and persists the key value pair
 */
func (s *RedisStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	session, sk, err := s.getSession(ctx, sessionID, userID)

	if err != nil {
		return Booking{}, err
	}

	s.rdb.Persist(ctx, sk)
	s.rdb.Persist(ctx, sessionKey(sessionID))

	session.Status = "confirmed"
	data := Booking{
		ID:      session.ID,
		MovieID: session.MovieID,
		UserID:  session.UserID,
		SeatID:  session.SeatID,
		Status:  session.Status,
	}

	val, _ := json.Marshal(data)
	s.rdb.Set(ctx, sk, val, 0)

	return session, nil
}

/*
* Parses Session from String JSON to Go Struct
 */
func parseSession(val string) (Booking, error) {
	var data Booking
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return Booking{}, err
	}
	return Booking{
		ID:      data.ID,
		MovieID: data.MovieID,
		SeatID:  data.SeatID,
		UserID:  data.UserID,
		Status:  data.Status,
	}, nil
}

/*
* Gets session from session ID and user ID
 */
func (s *RedisStore) getSession(ctx context.Context, sessionID string, userID string) (Booking, string, error) {
	sk, err := s.rdb.Get(ctx, sessionKey(sessionID)).Result()

	if err != nil {
		return Booking{}, "", err
	}

	val, err := s.rdb.Get(ctx, sk).Result()
	if err != nil {
		return Booking{}, "", err
	}

	session, err := parseSession(val)
	if err != nil {
		return Booking{}, "", err
	}

	return session, sk, nil
}

/*
* Deletes seat:session key value pair and sessionId:seat key value pair from Redis
* and frees up the seat for others
 */
func (s *RedisStore) Release(ctx context.Context, sessionID string, userID string) error {
	_, sk, err := s.getSession(ctx, sessionID, userID)

	if err != nil {
		return err
	}

	s.rdb.Del(ctx, sk, sessionKey(sessionID))
	return nil
}
