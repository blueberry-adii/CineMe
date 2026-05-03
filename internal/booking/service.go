package booking

import "context"

/*
* Service is required by HTTP Handler and Depends on BookingStore Interface
 */
type Service struct {
	store BookingStore
}

/*
* Constructor
 */
func NewService(store BookingStore) *Service {
	return &Service{
		store,
	}
}

/*
* Calls BookingStore Book to book a booking
 */
func (s *Service) Book(b Booking) (Booking, error) {
	return s.store.Book(b)
}

/*
* Calls BookingStore ListBookings and returns all the bookings for a
* given movieId
 */
func (s *Service) ListBookings(movieId string) []Booking {
	return s.store.ListBookings(movieId)
}

/*
* Calls BookingStore Confirm and Confirms the Booking for a given session and user
 */
func (s *Service) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	return s.store.Confirm(ctx, sessionID, userID)
}

/*
* Cancels the session and Releases seat held by the user
 */
func (s *Service) Release(ctx context.Context, sessionID string, userID string) error {
	return s.store.Release(ctx, sessionID, userID)
}
