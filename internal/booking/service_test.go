package booking

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/blueberry-adii/CineMe/internal/redis"
	"github.com/google/uuid"
)

/*
* Test to check for race conditions and atomicity
* memory_store failed - due to lack of atomicity in concurrent operations
* concurrent_store passed - due to use of mutex locks on read and write operations
* redis_store passed - due to redis being single threaded, writes are made only if session for
* given seat doesnt exist in redis.
 */
func TestConcurrentBooking(t *testing.T) {
	store := NewRedisStore(redis.NewRedisClient("localhost:6379"))
	svc := NewService(store)

	const goroutines = 100_000

	var (
		successes atomic.Int64
		failures  atomic.Int64
		wg        sync.WaitGroup
	)

	wg.Add(goroutines)
	for i := range goroutines {
		go func(userNum int) {
			defer wg.Done()
			_, err := svc.Book(Booking{
				MovieID: "screen-1",
				SeatID:  "A1",
				UserID:  uuid.New().String(),
			})

			if err == nil {
				successes.Add(1)
			} else {
				failures.Add(1)
			}
		}(i)
	}

	wg.Wait()

	/*
	* Only One User can succeed in booking 1 seat, all others fail - Test Pass
	* Otherwise - Test Fail
	 */
	if got := successes.Load(); got != 1 {
		t.Errorf("expected exactly 1 success, got %d", got)
	}
	if got := failures.Load(); got != int64(goroutines-1) {
		t.Errorf("expected exactly %d failures, got %d", goroutines-1, got)
	}
}
