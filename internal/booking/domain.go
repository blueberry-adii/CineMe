package booking

import (
	"context"
	"errors"
	"time"
)

/*
* Global Error Definitions
 */
var (
	ErrSeatAlreadyBooked = errors.New("seat is already taken")
)

/*
* A Booking contains Movie, Seat and the User, with a status
* and ExpiryDateAndTime if status is hold
 */
type Booking struct {
	ID        string
	MovieID   string
	SeatID    string
	UserID    string
	Status    string
	ExpiresAt time.Time
}

/*
* Interface for all Store structs
 */
type BookingStore interface {
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking

	Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
	Release(ctx context.Context, sessionID string, userID string) error
}
