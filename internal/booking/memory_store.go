package booking

/*
* Store 1:
* Memory store was created to test multiple user bookings, while using
* in memory bookings map
* Result:
* Failure, Race Conditions as multiple user overrid the bookings as the operations
* were not Atomic
 */
type MemoryStore struct {
	bookings map[string]Booking
}

/*
* Constructor
 */
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		map[string]Booking{},
	}
}

/*
* Adds a booking to the map if booking doesnt already exist
 */
func (s *MemoryStore) Book(b Booking) error {
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b
	return nil
}

/*
* Lists all the bookings in the bookings map
 */
func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var result []Booking

	for _, b := range s.bookings {
		if movieID == b.MovieID {
			result = append(result, b)
		}
	}

	return result
}
