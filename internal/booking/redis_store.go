package booking

import goredis "github.com/redis/go-redis/v9"

type RedisStore struct {
	rdb      *goredis.Client
	bookings map[string]Booking
}

func NewRedisStore(rdb *goredis.Client) *RedisStore {
	return &RedisStore{
		rdb,
		map[string]Booking{},
	}
}

func (s *RedisStore) Book(b Booking) error {
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b
	return nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	var result []Booking

	for _, b := range s.bookings {
		if movieID == b.MovieID {
			result = append(result, b)
		}
	}

	return result
}

// func (s *RedisStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {

// }
// func (s *RedisStore) Release(ctx context.Context, sessionID string, userID string) error {

// }
