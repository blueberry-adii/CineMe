package booking

type MemoryStore struct {
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) error {
	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b
	return nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var result []Booking

	for _, b := range s.bookings {
		if movieID == b.MovieID {
			result = append(result, b)
		}
	}

	return result
}

// func (s *MemoryStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {

// }
// func (s *MemoryStore) Release(ctx context.Context, sessionID string, userID string) error {

// }
