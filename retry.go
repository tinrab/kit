package util

import (
	"time"
)

type RetryFunc func(int) error

// Retry retries calling function f n-times.
// It returns an error if none of the tries succeeds.
func Retry(n int, f RetryFunc) (err error) {
	for i := 0; i < n; i++ {
		err = f(i)
		if err == nil {
			return nil
		}
	}
	return err
}

// RetrySleep retries calling function f n-times and sleeps for d after each call.
// It returns an error if none of the tries succeeds.
func RetrySleep(n int, d time.Duration, f RetryFunc) (err error) {
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
func RetryForever(f RetryFunc) {
	for i := 0; ; i++ {
		err := f(i)
		if err == nil {
			return
		}
	}
}

// RetryForeverSleep keeps trying to call function f until it succeeds, and sleeps after each failure.
func RetryForeverSleep(d time.Duration, f RetryFunc) {
	for i := 0; ; i++ {
		err := f(i)
		if err == nil {
			return
		}
		time.Sleep(d)
	}
}
