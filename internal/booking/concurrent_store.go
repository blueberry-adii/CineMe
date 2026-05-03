package booking

import "sync"

/*
* Store 2:
* Concurrent store was created to fix memory store's race conditions
* Implements Mutex Locks on Reads and Writes to prevent overriding
* Result:
* Success, no 2 users can write at the same time
 */
type ConcurrentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

/*
* Constructor
 */
func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
	}
}

/*
* Adds a booking to the map if booking doesnt already exist
 */
func (s *ConcurrentStore) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b
	return nil
}

/*
* Lists all the bookings in the bookings map
 */
func (s *ConcurrentStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock()

	var result []Booking
	for _, b := range s.bookings {
		if movieID == b.MovieID {
			result = append(result, b)
		}
	}

	return result
}
