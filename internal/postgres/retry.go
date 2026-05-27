package postgres

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/lib/pq"
)

const (
	retryAttempts  = 3
	retryBaseDelay = 50 * time.Millisecond
)

func WithRetry(ctx context.Context, fn func(context.Context) error) error {
	delay := retryBaseDelay
	for attempt := 1; attempt <= retryAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}
		if !ShouldRetry(err) {
			return err
		}
		if attempt == retryAttempts {
			return err
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
		delay *= 2
	}

	return nil
}

func ShouldRetry(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "40001", // serialization_failure
			"40P01": // deadlock_detected
			return true
		}

		// Connection-level failures
		if pqErr.Code.Class() == "08" {
			return true
		}
	}

	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
