package retry

import (
	"time"
)

type Func func(int) error

// Retry retries calling function f n-times.
// It returns an error if none of the tries succeeds.
func Retry(n int, f Func) (err error) {
	for i := 0; i < n; i++ {
		err = f(i)
		if err == nil {
			return nil
		}
	}
	return err
}

// Sleep retries calling function f n-times and sleeps for d after each call.
// It returns an error if none of the tries succeeds.
func Sleep(n int, d time.Duration, f Func) (err error) {
	for i := 0; i < n; i++ {
		err = f(i)
		if err == nil {
			return nil
		}
		time.Sleep(d)
	}
	return err
}

// RetryForever keeps trying to call function f until it succeeds.
func Forever(f Func) {
	for i := 0; ; i++ {
		err := f(i)
		if err == nil {
			return
		}
	}
}

// RetryForeverSleep keeps trying to call function f until it succeeds, and sleeps after each failure.
func ForeverSleep(d time.Duration, f Func) {
	for i := 0; ; i++ {
		err := f(i)
		if err == nil {
			return
		}
		time.Sleep(d)
	}
}
